package peer

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func LoadImage(ctx context.Context, baseURL string, imageContent io.Reader) error {
	targetURL := fmt.Sprintf("%s%s", baseURL, URLImageLoad)
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		written, _ := io.Copy(pipeWriter, imageContent)
		logger().WithName("client").Info("image content pipeWriter finished", "url", targetURL, "size", written)
	}()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, pipeReader)
	if err != nil {
		return errors.Wrap(err, "new request for load images")
	}

	response, err := http.DefaultClient.Do(request)
	defer func() {
		if err := response.Body.Close(); err != nil {
			logger().WithName("client").Error(err, "close response body")
		}
	}()

	return errors.Wrap(err, "send request for load images")
}
