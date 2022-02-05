package agent

import (
	"context"
	"io"
	"net/http"

	containerruntime "github.com/STRRL/kubectl-image/pkg/container/runtime"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

// URLImageLoad is the default HTTP endpoint for loading the container image.
// TODO: move HTTP server from cmd/kubectl-image-agent to here.
const URLImageLoad = "/image/load"

// Application is the core component for command kubectl-image-agent.
type Application struct {
	listenAddress string
	remote        containerruntime.Remote
	logger        logr.Logger
	httpServer    *http.Server
	running       *atomic.Bool
}

// NewApplication is the constructor for Application.
func NewApplication(listenAddress string, remote containerruntime.Remote, logger logr.Logger) *Application {
	result := &Application{
		listenAddress: listenAddress,
		remote:        remote,
		logger:        logger,
		running:       atomic.NewBool(false),
	}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
	})
	http.DefaultServeMux.HandleFunc(URLImageLoad, func(responseWriter http.ResponseWriter, r *http.Request) {
		err := result.forwardToDockerImageImport(r.Body)
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			_, _ = responseWriter.Write([]byte(err.Error()))

			return
		}
		responseWriter.WriteHeader(http.StatusOK)
	})

	result.httpServer = &http.Server{
		Addr:    result.listenAddress,
		Handler: serveMux,
	}

	return result
}

// Start would start the application with given context.
func (it *Application) Start(ctx context.Context) error {
	if !it.running.CAS(false, true) {
		return errors.New("already running")
	}

	defer func() {
		it.running.Store(false)
	}()

	go func() {
		err := it.httpServer.ListenAndServe()
		if err != nil {
			// TODO: resolve this error
			it.logger.Error(err, "start http server")
		}
	}()

	<-ctx.Done()

	if err := it.httpServer.Close(); err != nil {
		it.logger.Error(err, "close http server")
	}

	return nil
}

func (it *Application) forwardToDockerImageImport(body io.ReadCloser) error {
	err := it.remote.LoadImage(body)

	return errors.Wrapf(err, "forward to docker image import")
}
