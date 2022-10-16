package model

type Response struct {
	Services []Service `json:"services"`
	Nodes    Nodes     `json:"nodes"`
}

type Service struct {
	Name   string `json:"name"`
	Node   string `json:"node"`
	Status string `json:"status"`
}

type Nodes map[string]NodeStatus

type NodeStatus struct {
	Status string `json:"status"`
}
