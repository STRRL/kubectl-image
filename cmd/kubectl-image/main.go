package main

import (
	"github.com/STRRL/kubectl-push/pkg/cmd"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	defer func() {
		_ = zap.L().Sync()
	}()

	flags := pflag.NewFlagSet("kubectl-image", pflag.ExitOnError)
	pflag.CommandLine = flags

	zap.L().Info("Starting kubectl-push")

	err := cmd.NewCmdPush().Execute()
	if err != nil {
		zap.L().Error("failed to push image", zap.Error(err))
	}
}
