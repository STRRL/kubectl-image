# kubectl-push

[![Latest Docker Image](https://github.com/STRRL/kubectl-push/actions/workflows/latest-docker-image.yml/badge.svg)](https://github.com/STRRL/kubectl-push/actions/workflows/latest-docker-image.yml)
[![golangci-lint](https://github.com/STRRL/kubectl-push/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/STRRL/kubectl-push/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/STRRL/kubectl-push)](https://goreportcard.com/report/github.com/STRRL/kubectl-push)

docker push but for kubernetes

(WIP)

## Overview

I was always bothering with one problem: when I develop or debug chaos mesh, I need deliver the latest modified image to the target kubernetes cluster. If the kubernetes runs on `minikube` or `kind` with single node, I could use `kind load docker-image`, `minikube image load`, or build with`eval minikube docker-env`

## Feature

Deliver container images to kubernetes cluster simply, please do not use it in production env.

## Roadmap

- [x] build `kubectl-push-peer` that forwards HTTP request content to `docker image load`
- [x] play `kubectl-push-peer` with curl
- [x] build `kubectl-push`, automatically create Pod on Node, then send certain requests, in one command
- [ ] terminal UI progress, silent mode
- [x] support kubernetes container runtime: docker
- [x] support local container runtime: docker
- [ ] contributes to krew
- [ ] reduce duplicated layers
- [ ] support kubernetes container runtime: containerd
- [ ] support local container runtime: containerd
- [ ] support kubernetes container runtime: cri-o
- [ ] support local container runtime:  cri-o
- [ ] p2p transmission

## How it works

https://github.com/STRRL/kubectl-push/wiki/How-it-works%3F
