package importer

import (
	"io"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var logger = zapr.NewLogger(zap.L()).WithName("loader")

type ImageLoader interface {
	LoadImage(content io.ReadCloser) error
}
