package peer

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func LoadImage(ctx context.Context, baseUrl string, imageContent io.Reader) error {
	targetUrl := fmt.Sprintf("%s%s", baseUrl, UrlImageLoad)
	rx, tx := io.Pipe()

	go func() {
		written, _ := io.Copy(tx, imageContent)
		logger().WithName("client").Info("image content tx finished", "url", targetUrl, "size", written)
	}()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, targetUrl, rx)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(request)
	return err
}
