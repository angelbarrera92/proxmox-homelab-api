# Proxmox Homelab API

[![GitHub Super-Linter](https://github.com/angelbarrera92/proxmox-homelab-api/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)

## Usage

### CLI

```bash
$ ./proxmox-homelab-api -h
Usage: proxmox-homelab-api [options]
Options:
  -h, --help            Print this help
  -c, --config          Path to config file
  -v, --version         Print version information
```

## Build

```bash
go build -o proxmox-homelab-api cmd/proxmox-homelab-api/main.go
# or
make build
```

## Run

The subcommands are:

```bash
./proxmox-homelab-api -c configs/demo.yaml
```

### EndPoints

- `GET /services`: Returns the list of services and nodes with the status of each one.
- `POST /services/{service}/start`: Starts a service by name. It uses the proxmox API to start the VM.
- `POST /services/{service}/stop`: Stop a service by name. It uses the proxmox API to stop the VM.
- `POST /nodes/{node}/start`: Starts a node by name. It uses wake on lan to start the node.
- `POST /nodes/{node}/stop`: Stop a node by name. It uses the proxmox API to stop the node.

## Development

Use the `Makefile` to build and lint the code.

Requirements:

- `make`
- [`Docker`](https://docs.docker.com/get-docker/)

```bash
make clean lint build
```

## License

[MIT](LICENSE)
