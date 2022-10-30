package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/config"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/handlers"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/model"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/probes"
	"github.com/gorilla/mux"
)

var (
	proxmoxAPIHandler handlers.ProxmoxHomelabAPI
	data              model.Response
	logLevel          string
	version           = "dev"
	commit            = "none"
	date              = "unknown"
)

func printHelp() {
	fmt.Println("Usage: proxmox-homelab-api [options]")
	fmt.Println("Options:")
	fmt.Println("  -h, --help\t\tPrint this help")
	fmt.Println("  -c, --config\t\tPath to config file")
	fmt.Println("  -v, --version\t\tPrint version information")
}

func main() {

	// Print help if no arguments are provided
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	// Parse arguments
	var configPath string
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("proxmox-homelab-api %s, commit %s, built at %s\n", version, commit, date)
			os.Exit(0)
		case "-c", "--config":
			if i+1 < len(os.Args) {
				configPath = os.Args[i+1]
				i++
			} else {
				fmt.Println("Error: -c or --config requires a path to a config file")
				os.Exit(1)
			}
		default:
			fmt.Printf("Error: unknown argument %s\n", os.Args[i])
			os.Exit(1)
		}
	}

	// Load config
	cfg, err := config.Parse(configPath)
	if err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	// Init response
	data = initResponse(*cfg)

	// Init handlers
	proxmoxAPIHandler = handlers.ProxmoxHomelabAPI{
		Data:   &data,
		Config: *cfg,
	}

	// Configure logging
	logLevel = cfg.LogLevel
	if logLevel == "" {
		logLevel = "info"
	}

	// Run in background
	startAsyncProbes(*cfg, &data)

	// Create router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/services", proxmoxAPIHandler.Services)
	router.HandleFunc("/services/{service}/start", proxmoxAPIHandler.StartService)
	router.HandleFunc("/services/{service}/stop", proxmoxAPIHandler.StopService)
	router.HandleFunc("/nodes/{node}/start", proxmoxAPIHandler.StartNode)
	router.HandleFunc("/nodes/{node}/stop", proxmoxAPIHandler.StopNode)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           router,
	}

	log.Println("Starting server at", addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func initResponse(cfg config.Config) (d model.Response) {
	d = model.Response{
		Services: make([]model.Service, 0),
		Nodes:    model.Nodes{},
	}

	for _, node := range cfg.Nodes {
		d.Nodes[node.Name] = model.NodeStatus{
			Description: node.Description,
			EndPoint:    fmt.Sprintf("%s://%s:%d", node.Proxmox.Schema, node.Host, node.Proxmox.Port),
			Status:      "unknown",
			Probes:      len(node.Probes),
		}
	}

	for _, service := range cfg.Services {
		d.Services = append(d.Services, model.Service{
			Name:        service.Name,
			Description: service.Description,
			Icon:        service.Icon,
			Node:        service.Node,
			Status:      "unknown",
			Probes:      len(service.Probes),
			EndPoint:    fmt.Sprintf("%s://%s:%d", service.Schema, service.Host, service.Port),
		})
	}

	return
}

func startAsyncProbes(cfg config.Config, data *model.Response) {
	go probes.NodeProbes(cfg.Nodes, data)
	go probes.ServiceProbes(cfg.Services, data)
}
