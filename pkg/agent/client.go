package agent

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

// Client is the interface to interact with each kubectl-image-agent.
type Client interface {
	// Health execute the health check for target kubectl-image-agent.
	Health(ctx context.Context) (bool, error)

	// LoadImage loads the image into the nodes which running kubectl-image-agent.
	LoadImage(ctx context.Context, imageContent io.Reader) error

	// ListImage return all the existed container images
	ListImage(ctx context.Context) ([]ContainerImage, error)
}

// HTTPClient is the default implementation of the Client interface.
type HTTPClient struct {
	address string
	logger  logr.Logger
}

// NewHTTPClient is the constructor for the HTTPClient.
func NewHTTPClient(address string, logger logr.Logger) *HTTPClient {
	return &HTTPClient{address: address, logger: logger}
}

// Health implements the agent.Client interface.
func (it *HTTPClient) Health(ctx context.Context) (bool, error) {
	targetURL := fmt.Sprintf("%s/%s", it.address, URLHealth)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return false, errors.Wrap(err, "create HTTP request")
	}

	response, err := http.DefaultClient.Do(request)
	defer func() {
		_ = response.Body.Close()
	}()

	if err != nil {
		return false, errors.Wrapf(err, "probe health for %s", targetURL)
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return true, nil
	}

	return false, nil
}

// LoadImage implements the agent.Client interface.
func (it *HTTPClient) LoadImage(ctx context.Context, imageContent io.Reader) error {
	targetURL := fmt.Sprintf("%s%s", it.address, URLImageLoad)
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		written, _ := io.Copy(pipeWriter, imageContent)
		it.logger.Info("image content pipeWriter finished", "url", targetURL, "size", written)
	}()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, pipeReader)
	if err != nil {
		return errors.Wrap(err, "new request for load images")
	}

	response, err := http.DefaultClient.Do(request)
	defer func() {
		if err := response.Body.Close(); err != nil {
			it.logger.Error(err, "close response body")
		}
	}()

	return errors.Wrap(err, "send request for load images")
}
