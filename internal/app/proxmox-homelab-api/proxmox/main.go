package proxmox

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/Telmate/proxmox-api-go/proxmox"
)

type ClientConfig struct {
	Schema    string
	Host      string
	Port      int
	Username  string
	Password  string
	VerifySSL bool
}

func NewClient(config ClientConfig) (*proxmox.Client, error) {
	proxmoxAPIURL := fmt.Sprintf("%s://%s:%d/api2/json", config.Schema, config.Host, config.Port)
	insecure := !config.VerifySSL
	proxyString := ""
	httpHeaders := ""
	taskTimeout := 300
	tlsConf := &tls.Config{InsecureSkipVerify: !config.VerifySSL} // nolint:gosec
	if !insecure {
		tlsConf = nil
	}

	c, err := proxmox.NewClient(proxmoxAPIURL, nil, httpHeaders, tlsConf, proxyString, taskTimeout)
	if err != nil {
		log.Fatalf("error creating proxmox client: %s", err)
		return nil, err
	}

	err = c.Login(config.Username, config.Password, "")
	if err != nil {
		log.Fatalf("error logging in: %s", err)
		return nil, err
	}

	return c, nil
}
