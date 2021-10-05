package main

import (
	"io"
	"net/http"

	importer "github.com/STRRL/kubectl-push/pkg/loader"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var logger = zapr.NewLogger(zap.L()).WithName("main")

func main() {
	http.DefaultServeMux.HandleFunc("/image/load", func(rw http.ResponseWriter, r *http.Request) {
		err := forwardToDockerImageImport(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(":28375", http.DefaultServeMux)
	if err != nil {
		logger.Error(err, "failed to start http server")
	}
}

func forwardToDockerImageImport(content io.ReadCloser) error {
	cr := importer.DockerImageLoader{}
	return cr.LoadImage(content)
}
