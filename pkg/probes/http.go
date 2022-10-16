package probes

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type HTTPProber struct {
	Schema         string
	Host           string
	Port           int
	Path           string
	ExpectedStatus int
}

func (p HTTPProber) Probe() (bool, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}
	probeEndpoint := fmt.Sprintf("%s://%s:%d%s", p.Schema, p.Host, p.Port, p.Path)
	resp, err := client.Get(probeEndpoint)
	if err != nil {
		return false, fmt.Errorf("error probing http: %w", err)
	}
	if resp.StatusCode == p.ExpectedStatus {
		return true, nil
	}
	return false, nil
}
