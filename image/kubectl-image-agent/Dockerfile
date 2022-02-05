# syntax=docker/dockerfile:experimental

FROM bitnami/golang:1.17.2 AS build-env
WORKDIR /
COPY . /kubectl-image
ENV GO111MODULE on
RUN cd /kubectl-image && make

FROM debian:stable-slim
LABEL org.opencontainers.image.source=https://github.com/STRRL/kubectl-image
COPY --from=docker:dind /usr/local/bin/docker /usr/local/bin/
COPY --from=build-env /kubectl-image/bin/kubectl-image-agent /usr/local/bin/kubectl-image-agent
CMD kubectl-image-agent
