package probes

import (
	"time"

	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/config"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/model"
)

func NodeProbes(nodes []config.Node, data *model.Response) {
	for {

		// Prepare async probes
		asyncProbes := make([]asyncProbe, 0)
		for _, node := range nodes {
			name := node.Name
			host := node.Host
			asyncProbes = append(asyncProbes, createAsyncProbes(name, host, node.Probes)...)
		}

		// Run async probes
		probeResult := runProbers(asyncProbes)

		// Update data
		data.Nodes = model.Nodes{}
		for _, node := range nodes {
			var nodeStatus model.NodeStatus
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
