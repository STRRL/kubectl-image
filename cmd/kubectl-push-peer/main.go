package main

import (
	"context"

	containerruntime "github.com/STRRL/kubectl-push/pkg/container/runtime"
	"github.com/STRRL/kubectl-push/pkg/peer"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func main() {
	logger := zapr.NewLogger(zap.L()).WithName("main")

	application := peer.NewApplication("0.0.0.0:28375", &containerruntime.DockerCli{}, logger)

	if err := application.Start(context.TODO()); err != nil {
		logger.Error(err, "failed to start http server")
	}
}
