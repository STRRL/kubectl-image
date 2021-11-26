package main

import (
	"github.com/STRRL/kubectl-push/pkg/peer"
	"io"
	"net/http"

	containerruntime "github.com/STRRL/kubectl-push/pkg/container/runtime"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var logger = zapr.NewLogger(zap.L()).WithName("main")

func main() {
	http.DefaultServeMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.DefaultServeMux.HandleFunc(peer.UrlImageLoad, func(rw http.ResponseWriter, r *http.Request) {
		err := forwardToDockerImageImport(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe("0.0.0.0:28375", http.DefaultServeMux)
	if err != nil {
		logger.Error(err, "failed to start http server")
	}
}

func forwardToDockerImageImport(content io.ReadCloser) error {
	cr := &containerruntime.Docker{}
	return cr.LoadImage(content)
}
