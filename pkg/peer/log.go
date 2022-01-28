package peer

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func logger() logr.Logger {
	return zapr.NewLogger(zap.L()).WithName("peer")
}
