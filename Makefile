.PHONY: all
all: clean binary

.PHONY: binary
binary: bin/kubectl-push bin/kubectl-push-peer

bin/kubectl-push:
	go build -o ./bin/kubectl-push ./cmd/kubectl-push/main.go

bin/kubectl-push-peer:
	go build -o ./bin/kubectl-push-peer ./cmd/kubectl-push-peer/main.go

.PHONY: image
image: image/kubectl-push-peer

.PHONY: image/kubectl-push-peer
image/kubectl-push-peer: bin/kubectl-push-peer
	DOCKER_BUILDKIT=0 docker build -t ghcr.io/strrl/kubectl-push-peer:latest ./image/kubectl-push-peer

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: check
check:
	golangci-lint run
