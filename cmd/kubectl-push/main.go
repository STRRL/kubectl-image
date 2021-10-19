package main

import (
	"github.com/STRRL/kubectl-push/pkg/cmd"
	"go.uber.org/zap"
)

func init() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	zap.L().Info("Starting kubectl-push")

	cmd.NewCmdPush().Execute()
}
