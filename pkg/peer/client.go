package peer

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func LoadImage(ctx context.Context, baseUrl string, imageContent io.Reader) error {
	_, err := http.Post(fmt.Sprintf("%s/image/load", baseUrl), "application/octet-stream", imageContent)
	return err
}
