.PHONY: all
all: clean binary

.PHONY: binary
binary: bin/kubectl-image bin/kubectl-image-agent

bin/kubectl-image:
	go build -o ./bin/kubectl-image ./cmd/kubectl-image/main.go

bin/kubectl-image-agent:
	go build -o ./bin/kubectl-image-agent ./cmd/kubectl-image-agent/main.go

.PHONY: image
image: image/kubectl-image-agent

.PHONY: image/kubectl-image-agent
image/kubectl-image-agent: bin/kubectl-image-agent
	DOCKER_BUILDKIT=0 docker build -t ghcr.io/strrl/kubectl-image-agent:latest -f image/kubectl-image-agent/Dockerfile .

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: check
check:
	golangci-lint run --fix
