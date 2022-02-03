package runtime

import (
	"context"
	"github.com/docker/docker/client"
	"io"
	"time"
)

var _ Remote = (*Docker)(nil)
var _ Local = (*Docker)(nil)

type Docker struct {
	imageAPIClient client.ImageAPIClient
	timeout        time.Duration
}

func (it *Docker) LoadImage(content io.ReadCloser) error {
	ctx, cancelFunc := context.WithDeadline(context.TODO(), time.Now().Add(it.timeout))
	defer cancelFunc()

	_, err := it.imageAPIClient.ImageLoad(ctx, content, false)
	return err
}

func (it *Docker) ImageExist(imageName string) (bool, error) {
	ctx, cancelFunc := context.WithDeadline(context.TODO(), time.Now().Add(it.timeout))
	defer cancelFunc()

	_, _, err := it.imageAPIClient.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (it *Docker) ImageSave(imageName string, content io.Writer) error {
	ctx, cancelFunc := context.WithDeadline(context.TODO(), time.Now().Add(it.timeout))
	defer cancelFunc()

	reader, err := it.imageAPIClient.ImageSave(ctx, []string{imageName})
	if err != nil {
		return err
	}

	_, err = io.Copy(content, reader)
	return err
}
