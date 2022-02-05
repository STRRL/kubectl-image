# kubectl-image

[![Latest Docker Image](https://github.com/STRRL/kubectl-image/actions/workflows/latest-docker-image.yml/badge.svg)](https://github.com/STRRL/kubectl-push/actions/workflows/latest-docker-image.yml)
[![golangci-lint](https://github.com/STRRL/kubectl-image/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/STRRL/kubectl-push/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/STRRL/kubectl-image)](https://goreportcard.com/report/github.com/STRRL/kubectl-push)

`docker image` but for kubernetes

(WIP)

## Overview

Kubernetes (nearly) does not care about what images exist on the node, the only thing relates to "managing" image is image garbage collection.

I think it's not so convenient to cluster admin using another certain tools to management the images.

And another thing is I was always bothering with one problem: when I develop or debug chaos mesh, I need deliver the latest modified image to the target kubernetes cluster. If the kubernetes runs on `minikube` or `kind` with single node, I could use `kind load docker-image`, `minikube image load`, or build with`eval minikube docker-env`

## Feature

- List images on each Kubernetes Node
- load image from local image
- Deliver container images to kubernetes cluster simply, please do not use it in production env.

## Roadmap

> Out-of-date: This roadmap need updates.

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

https://github.com/STRRL/kubectl-image/wiki/How-it-works%3F
