package provisioner

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func getLogger() logr.Logger {
	return zapr.NewLogger(zap.L()).WithName("push")
}
