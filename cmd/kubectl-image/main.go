package main

import (
	stdlog "log"

	"github.com/STRRL/kubectl-image/pkg/kubectlimage/cmd"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		stdlog.Fatal(err)
	}

	logger := zapr.NewLogger(zapLogger)

	err = cmd.NewRootCommand().Execute()
	if err != nil {
		logger.Error(err, "kubectl image")
	}
}
