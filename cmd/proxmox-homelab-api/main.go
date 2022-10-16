package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/config"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/handlers"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/model"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/probes"
	"github.com/gorilla/mux"
)

var (
	proxmoxApiHandler handlers.ProxmoxHomelabApi
	data              model.Response
	logLevel          string
)

func main() {

	// Load config
	cfg, err := config.Parse("/home/angel/personal/proxmox-homelab-api/configs/demo.yaml")
	if err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	// Init response
	data = initResponse(*cfg)

	// Init handlers
	proxmoxApiHandler = handlers.ProxmoxHomelabApi{
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
	router.HandleFunc("/services", proxmoxApiHandler.Services)
	router.HandleFunc("/services/{service}/start", proxmoxApiHandler.StartService)
	router.HandleFunc("/services/{service}/stop", proxmoxApiHandler.StopService)
	router.HandleFunc("/nodes/{node}/start", proxmoxApiHandler.StartNode)
	router.HandleFunc("/nodes/{node}/stop", proxmoxApiHandler.StopNode)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Println("Starting server at", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func initResponse(cfg config.Config) (d model.Response) {
	d = model.Response{
		Services: make([]model.Service, 0),
		Nodes:    model.Nodes{},
	}

	for _, node := range cfg.Nodes {
		d.Nodes[node.Name] = model.NodeStatus{
			Status: "unknown",
		}
	}

	for _, service := range cfg.Services {
		d.Services = append(d.Services, model.Service{
			Name:   service.Name,
			Node:   service.Node,
			Status: "unknown",
		})
	}

	return
}

func startAsyncProbes(cfg config.Config, data *model.Response) {
	go probes.NodeProbes(cfg.Nodes, data)
	go probes.ServiceProbes(cfg.Services, data)
}
