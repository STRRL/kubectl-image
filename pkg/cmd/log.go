package cmd

import (
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var logger = zapr.NewLogger(zap.L()).WithName("push")
