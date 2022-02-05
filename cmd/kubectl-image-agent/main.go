package main

import (
	"context"

	"github.com/STRRL/kubectl-image/pkg/agent"
	containerruntime "github.com/STRRL/kubectl-image/pkg/agent/container/runtime"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func main() {
	logger := zapr.NewLogger(zap.L()).WithName("main")

	application := agent.NewApplication("0.0.0.0:28375", &containerruntime.DockerCli{}, logger)

	if err := application.Start(context.TODO()); err != nil {
		logger.Error(err, "failed to start http server")
	}
}
