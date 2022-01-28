package main

import (
	"io"
	"net/http"

	containerruntime "github.com/STRRL/kubectl-push/pkg/container/runtime"
	"github.com/STRRL/kubectl-push/pkg/peer"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	logger := zapr.NewLogger(zap.L()).WithName("main")

	http.DefaultServeMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.DefaultServeMux.HandleFunc(peer.URLImageLoad, func(responseWriter http.ResponseWriter, r *http.Request) {
		err := forwardToDockerImageImport(r.Body)
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			_, _ = responseWriter.Write([]byte(err.Error()))

			return
		}

		responseWriter.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe("0.0.0.0:28375", http.DefaultServeMux)
	if err != nil {
		logger.Error(err, "failed to start http server")
	}
}

func forwardToDockerImageImport(content io.ReadCloser) error {
	containerRuntime := &containerruntime.Docker{}
	err := containerRuntime.LoadImage(content)

	return errors.Wrapf(err, "forward content to container runtime image load")
}
