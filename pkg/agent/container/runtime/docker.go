package runtime

import (
	"context"
	"github.com/STRRL/kubectl-image/pkg/agent"
	"github.com/docker/docker/api/types"
	"io"
	"time"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

var (
	_ Remote = (*Docker)(nil)
	_ Local  = (*Docker)(nil)
)

// Docker is the another implementation for docker, the accessor of DockerCli.
type Docker struct {
	imageAPIClient client.ImageAPIClient
	timeout        time.Duration
}

// LoadImage implements the Remote.LoadImage.
func (it *Docker) LoadImage(ctx context.Context, content io.ReadCloser) error {
	_, err := it.imageAPIClient.ImageLoad(ctx, content, false)

	return errors.Wrap(err, "load image")
}

func (it *Docker) ListImages(ctx context.Context) ([]agent.ContainerImage, error) {
	list, err := it.imageAPIClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "list images")
	}
	var result []agent.ContainerImage
	for _, item := range list {
		result = append(result,
			agent.ContainerImage{
				Repository: item.RepoTags,
				Tag:        "",
				Digest:     "",
				ImageID:    "",
				Created:    time.Time{},
				Size:       0,
			})
	}
	return result
}

// ImageExist implements the Local.ImageExist.
func (it *Docker) ImageExist(imageName string) (bool, error) {
	ctx, cancelFunc := context.WithDeadline(context.TODO(), time.Now().Add(it.timeout))
	defer cancelFunc()

	_, _, err := it.imageAPIClient.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}

		return false, errors.Wrapf(err, "inspect image %s", imageName)
	}

	return true, nil
}

// ImageSave implements the Local.ImageSave.
func (it *Docker) ImageSave(imageName string, content io.Writer) error {
	ctx, cancelFunc := context.WithDeadline(context.TODO(), time.Now().Add(it.timeout))
	defer cancelFunc()

	reader, err := it.imageAPIClient.ImageSave(ctx, []string{imageName})
	if err != nil {
		return errors.Wrapf(err, "save image %s", imageName)
	}

	_, err = io.Copy(content, reader)

	return errors.Wrapf(err, "save image %s, copy image content", imageName)
}
