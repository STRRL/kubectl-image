FROM alpine:latest
RUN apk add --no-cache docker
COPY ./bin/kubectl-push-peer /kubectl-push-peer
CMD kubectl-push-peer
