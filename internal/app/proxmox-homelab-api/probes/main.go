package probes

import (
	"log"
	"sync"

	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/config"
	"github.com/angelbarrera92/proxmox-homelab-api/pkg/probes"
)

type asyncProbe struct {
	Name   string
	Prober probes.Prober
}

type asyncProbeResult struct {
	Name   string
	Result bool
	Error  error
}

type probeResults map[string]int

func createAsyncProbes(name, host string, hostProbes []config.Probe) (asyncProbes []asyncProbe) {
	asyncProbes = make([]asyncProbe, 0)

	for _, probe := range hostProbes {
		var prober probes.Prober
		if probe.Type == "port" {
			prober = probes.PortProber{
				Host: host,
				Port: probe.Port,
			}
		} else if probe.Type == "http" {
			prober = probes.HTTPProber{
				Schema:             probe.Schema,
				Host:               host,
				Port:               probe.Port,
				Path:               probe.Path,
				ExpectedStatus:     probe.ExpectedStatus,
				InsecureSkipVerify: probe.InsecureSkipVerify,
			}
		}

		asyncProbes = append(asyncProbes, asyncProbe{
			Name:   name,
			Prober: prober,
		})
	}

	return
}

func runProbers(asyncProbes []asyncProbe) probeResults {

	// Create a channel to receive the results
	results := make(chan asyncProbeResult, len(asyncProbes))

	// Create a wait group to wait for all the probes to finish
	wg := sync.WaitGroup{}

	// Loop over the probes and start a goroutine for each one
	for _, probe := range asyncProbes {
		wg.Add(1)
		go func(probe asyncProbe) {
			// Defer the wait group Done
			defer wg.Done()
			// Probe the target
			probeResult, probeError := probe.Prober.Probe()
			// Send the result to the results channel
			results <- asyncProbeResult{
				Name:   probe.Name,
				Result: probeResult,
				Error:  probeError,
			}
		}(probe)
	}

	// Wait for all the probes to finish
	wg.Wait()
	// Close the results channel
	close(results)

	pr := probeResults{}

	for result := range results {
		if result.Error != nil {
			log.Printf("error probing %s: %s", result.Name, result.Error)
		}
		if result.Result {
			pr[result.Name]++
		}
	}

	return pr
}
