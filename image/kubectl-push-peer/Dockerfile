# syntax=docker/dockerfile:experimental

FROM golang:1.16-alpine AS build-env
ENV GO111MODULE on
RUN go install github.com/STRRL/kubectl-push/cmd/kubectl-push-peer@latest

FROM alpine:latest
LABEL org.opencontainers.image.source https://github.com/STRRL/kubectl-push
RUN apk add --no-cache docker-cli
COPY --from=build-env /go/bin/kubectl-push-peer /usr/local/bin/kubectl-push-peer
CMD kubectl-push-peer
