.PHONY: all
all: clean binary

.PHONY: binary
binary: bin/kubectl-image bin/kubectl-push-peer

bin/kubectl-image:
	go build -o ./bin/kubectl-image ./cmd/kubectl-image/main.go

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
	golangci-lint run --fix
