package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host     string    `yaml:"host"`
	Port     int       `yaml:"port"`
	LogLevel string    `yaml:"log_level"`
	Nodes    []Node    `yaml:"nodes"`
	Services []Service `yaml:"services"`
}
type Wol struct {
	Broadcast string `yaml:"broadcast"`
}

type Probe struct {
	Type               string `yaml:"type"`
	Port               int    `yaml:"port"`
	Schema             string `yaml:"schema,omitempty"`
	Path               string `yaml:"path,omitempty"`
	ExpectedStatus     int    `yaml:"expected_status,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify,omitempty"`
}

type Node struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Host        string  `yaml:"host"`
	Mac         string  `yaml:"mac"`
	Wol         Wol     `yaml:"wol"`
	Proxmox     Proxmox `yaml:"proxmox"`
	Probes      []Probe `yaml:"probes"`
}

type Proxmox struct {
	Schema    string `yaml:"schema"`
	VerifySSL bool   `yaml:"verifySSL"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type Service struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Node        string  `yaml:"node"`
	Host        string  `yaml:"host"`
	Port        int     `yaml:"port"`
	Schema      string  `yaml:"schema"`
	Probes      []Probe `yaml:"probes"`
}

func Parse(filepath string) (*Config, error) {
	// Read the yaml file and unmarshal it into a Config struct
	var config Config
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}
	return &config, nil
}
