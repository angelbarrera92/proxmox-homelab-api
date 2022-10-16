# proxmox-homelab-api

## curl

```bash
$ curl -v localhost:8080/services | jq
$ curl -X POST localhost:8080/nodes/brix2807/start
$ curl -X POST localhost:8080/nodes/brix2807/stop
```