package handlers

import (
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/config"
	"github.com/angelbarrera92/proxmox-homelab-api/internal/app/proxmox-homelab-api/model"
)

type ProxmoxHomelabAPI struct {
	Data   *model.Response
	Config config.Config
}
