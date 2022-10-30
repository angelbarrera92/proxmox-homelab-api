package model

type Response struct {
	Services []Service `json:"services"`
	Nodes    Nodes     `json:"nodes"`
}

type Service struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	EndPoint    string `json:"endpoint"`
	Node        string `json:"node"`
	Probes      int    `json:"probes"`
	Status      string `json:"status"`
}

type Nodes map[string]NodeStatus

type NodeStatus struct {
	Description string `json:"description"`
	Status      string `json:"status"`
	EndPoint    string `json:"endpoint"`
	Probes      int    `json:"probes"`
}
