VERSION ?= dev
COMMIT ?= none
DATE ?= unknown

build: clean
	@docker run --rm -v $(shell mktemp -d):/cache -v $(shell pwd):/home/proxmox-homelab-api -u $(shell id -u):$(shell id -g) -w /home/proxmox-homelab-api -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 -e GOMODCACHE=/cache/mod -e GOCACHE=/cache/build cgr.dev/chainguard/go:1.19 build -o proxmox-homelab-api -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'"  cmd/proxmox-homelab-api/main.go

lint:
	@docker run --rm -e RUN_LOCAL=true -v $(shell pwd):/tmp/lint github/super-linter:v4

clean:
	@rm -rf proxmox-homelab-api
	@rm -rf super-linter.log

container: build
	@docker build -t proxmox-homelab-api:local -f build/container/Dockerfile .
