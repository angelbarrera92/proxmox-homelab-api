package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/proxmox"
	"github.com/gorilla/mux"
)

type ACTION string

const (
	STOP  ACTION = "stop"
	START ACTION = "start"
)

// Endpoint to get the status of all services
// nolint:typecheck
func (p ProxmoxHomelabAPI) Services(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p.Data)
}

// nolint:typecheck
func (p ProxmoxHomelabAPI) StartService(w http.ResponseWriter, r *http.Request) {
	p.manageService(START, w, r)
}

// nolint:typecheck
func (p ProxmoxHomelabAPI) StopService(w http.ResponseWriter, r *http.Request) {
	p.manageService(STOP, w, r)
}

// nolint:typecheck
func (p ProxmoxHomelabAPI) manageService(a ACTION, w http.ResponseWriter, r *http.Request) {
	// Check method, only POST is allowed
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	service := mux.Vars(r)["service"]

	// Check user submitted the service
	if service == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the array of services contains the service requested
	for _, s := range p.Config.Services {
		if s.Name == service {
			// Check if the array of services in data contains the service requested
			for _, ds := range p.Data.Services {
				if ds.Name == service {
					if err := checkServiceStatusForAction(a, ds.Status); err != nil {
						w.WriteHeader(http.StatusConflict)
						return
					}
				}
			}

			// Get its node from the array of nodes
			for _, n := range p.Config.Nodes {
				if n.Name == s.Node {

					// Check if the node exists in the map of nodes in data
					if _, ok := p.Data.Nodes[n.Name]; !ok {
						w.WriteHeader(http.StatusNotFound)
						return
					}

					// Check the node status
					if checkNodeStatusForAction(a, p.Data.Nodes[n.Name].Status) != nil {
						w.WriteHeader(http.StatusConflict)
						return
					}

					// TODO: Potentially, if node is not running, start it
					// Wait for the node to be running
					// Start the service
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
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					err = shutdownOrStartVM(s.Name, a, c)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					w.WriteHeader(http.StatusAccepted)
				}
			}
		}
	}
}

func checkServiceStatusForAction(a ACTION, s string) (err error) {
	switch a {
	case STOP:
		if s != "ok" {
			return fmt.Errorf("service is not running")
		}
	case START:
		if s != "error" {
			return fmt.Errorf("service is already running")
		}
	}

	return nil
}

func checkNodeStatusForAction(a ACTION, s string) (err error) {
	switch a {
	case STOP:
		if s != "ok" {
			return fmt.Errorf("node is not running")
		}
	case START:
		if s != "ok" {
			return fmt.Errorf("node is not running")
		}
	}

	return nil
}

func shutdownOrStartVM(vmName string, a ACTION, c *proxmox.Client) (err error) {
	vmRef, err := c.GetVMRefByName(vmName)
	if err != nil {
		return err
	}

	switch a {
	case STOP:
		_, err = c.ShutdownVM(vmRef)
		return
	case START:
		_, err = c.StartVM(vmRef)
		return
	}

	return
}
