package probes

import (
	"fmt"
	"time"

	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/config"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/model"
)

func ServiceProbes(services []config.Service, data *model.Response) {
	for {

		// Prepare async probes
		asyncProbes := make([]asyncProbe, 0) // nolint:typecheck
		for _, service := range services {
			name := service.Name
			host := service.Host
			asyncProbes = append(asyncProbes, createAsyncProbes(name, host, service.Probes)...) // nolint:typecheck
		}

		// Run async probes
		probeResult := runProbers(asyncProbes) // nolint:typecheck

		// Update data
		data.Services = make([]model.Service, len(services))

		for index, service := range services {
			serviceStatus := model.Service{
				Name:        service.Name,
				Description: service.Description,
				Icon:        service.Icon,
				Node:        service.Node,
				Probes:      len(service.Probes),
				EndPoint:    fmt.Sprintf("%s://%s:%d", service.Schema, service.Host, service.Port),
			}
			if probeResult[service.Name] == len(service.Probes) {
				serviceStatus.Status = "ok"
			} else {
				serviceStatus.Status = "error"
			}

			data.Services[index] = serviceStatus
		}

		// TODO: make this configurable
		time.Sleep(5 * time.Second)
	}
}
