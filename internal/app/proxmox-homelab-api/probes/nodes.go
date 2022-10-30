package probes

import (
	"fmt"
	"time"

	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/config"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/model"
)

func NodeProbes(nodes []config.Node, data *model.Response) {
	for {

		// Prepare async probes
		asyncProbes := make([]asyncProbe, 0) // nolint:typecheck
		for _, node := range nodes {
			name := node.Name
			host := node.Host
			asyncProbes = append(asyncProbes, createAsyncProbes(name, host, node.Probes)...) // nolint:typecheck
		}

		// Run async probes
		probeResult := runProbers(asyncProbes) // nolint:typecheck

		// Update data
		data.Nodes = model.Nodes{}
		for _, node := range nodes {
			nodeStatus := model.NodeStatus{
				Description: node.Description,
				EndPoint:    fmt.Sprintf("%s://%s:%d", node.Proxmox.Schema, node.Host, node.Proxmox.Port),
				Probes:      len(node.Probes),
			}
			if probeResult[node.Name] == len(node.Probes) {
				nodeStatus.Status = "ok"
			} else {
				nodeStatus.Status = "error"
			}

			data.Nodes[node.Name] = nodeStatus
		}

		// TODO: make this configurable
		time.Sleep(5 * time.Second)
	}
}
