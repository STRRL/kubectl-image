bin/kubectl-push-peer:
	go build -o ./bin/kubectl-push-peer ./cmd/kubectl-push-peer/main.go

.PHONY: image/kubectl-push-peer
image/kubectl-push-peer: bin/kubectl-push-peer
	DOCKER_BUILDKIT=0 docker build -t ghcr.io/strrl/kubectl-push-peer:latest -f ./image/kubectl-push-peer/Dockerfile .
