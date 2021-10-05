# kubectl-push

docker push but for kubernetes

## Overview

I was always bothering with one problem: when I develop or debug chaos mesh, I need deliver the latest modified image to the target kubernetes cluster. If the kubernetes runs on `minikube` or `kind` with single node, I could use `kind load docker-image`, `minikube image load`, or build with`eval minikube docker-env`

## How it works

https://github.com/STRRL/kubectl-push/wiki/How-it-works%3F
