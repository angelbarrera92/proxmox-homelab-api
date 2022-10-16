module github.com/angelbarrera92/proxmox-homelab-api

go 1.17

require (
	github.com/Telmate/proxmox-api-go v0.0.0-20221015123156-ba07d5ebc42b
	github.com/gorilla/mux v1.8.0
	github.com/sabhiram/go-wol v0.0.0-20211224004021-c83b0c2f887d
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/Telmate/proxmox-api-go => github.com/angelbarrera92/proxmox-api-go v0.0.0-20221014150916-a37ad72fe1b0
