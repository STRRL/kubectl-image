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
	rx, tx := io.Pipe()

	go func() {
		written, _ := io.Copy(tx, imageContent)
		logger().WithName("client").Info("image content tx finished", "url", targetURL, "size", written)
	}()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, rx)
	if err != nil {
		return errors.Wrap(err, "new request for load images")
	}
	_, err = http.DefaultClient.Do(request)
	return errors.Wrap(err, "send request for load images")
}
