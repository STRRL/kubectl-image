# syntax=docker/dockerfile:experimental

FROM bitnami/golang:1.17.2 AS build-env
WORKDIR /
RUN git clone --depth=1 https://github.com/STRRL/kubectl-push
ENV GO111MODULE on
RUN cd /kubectl-push && make

FROM debian:stable-slim
LABEL org.opencontainers.image.source https://github.com/STRRL/kubectl-push
COPY --from=docker:dind /usr/local/bin/docker /usr/local/bin/
COPY --from=build-env /kubectl-push/bin/kubectl-push-peer /usr/local/bin/kubectl-push-peer
CMD kubectl-push-peer
