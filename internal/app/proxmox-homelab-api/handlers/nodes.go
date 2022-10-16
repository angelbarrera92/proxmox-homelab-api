package handlers

import (
	"log"
	"net/http"

	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/proxmox"
	"github.com/angelbarrera92/proxmox-homelab-api/pkg/wol"
	"github.com/gorilla/mux"
)

// StartNode starts a node: POST /nodes/{node}/start
func (p ProxmoxHomelabApi) StartNode(w http.ResponseWriter, r *http.Request) {
	// Check method, only POST is allowed
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	node := mux.Vars(r)["node"]

	// Check user submitted the node
	if node == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := p.Data.Nodes[node]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if node is already running
	if p.Data.Nodes[node].Status == "ok" {
		w.WriteHeader(http.StatusConflict)
		return
	}

	err := p.startNode(node)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

func (p ProxmoxHomelabApi) startNode(node string) error {
	log.Printf("Starting node %s", node)
	for _, n := range p.Config.Nodes {
		if n.Name == node {
			log.Printf("Sending magic packet to %s", n.Mac)
			err := wol.SendMagicPacket(n.Mac, n.Wol.Broadcast)
			if err != nil {
				log.Printf("Error starting node %s: %v", node, err)
				return err
			}
			break
		}
	}
	return nil
}

// StartNode shutdown a node: POST /nodes/{node}/stop
func (p ProxmoxHomelabApi) StopNode(w http.ResponseWriter, r *http.Request) {
	// Check method, only POST is allowed
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	node := mux.Vars(r)["node"]

	// Check user submitted the node
	if node == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := p.Data.Nodes[node]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if node is already running
	if p.Data.Nodes[node].Status != "ok" {
		w.WriteHeader(http.StatusConflict)
		return
	}
	err := p.stopNode(node)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (p ProxmoxHomelabApi) stopNode(node string) error {
	log.Printf("Stopping node %s", node)
	for _, n := range p.Config.Nodes {
		if n.Name == node {
			// Create a new client
			config := proxmox.ClientConfig{
				Schema:    n.Proxmox.Schema,
				Host:      n.Host,
				Port:      n.Proxmox.Port,
				Username:  n.Proxmox.Username,
				Password:  n.Proxmox.Password,
				VerifySSL: n.Proxmox.VerifySSL,
			}
			c, err := proxmox.NewClient(config)
			if err != nil {
				log.Fatalf("error creating proxmox client: %s", err)
				return err
			}
			_, err = c.ShutdownNode(node)
			if err != nil {
				log.Fatalf("error shutting down node: %s", err)
				return err
			}

			break
		}
	}
	return nil
}
