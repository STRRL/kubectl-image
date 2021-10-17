package runtime

import "io"

// Interface Local means that the container runtime is running at the "client side"/"kubectl side".
// Interface Local check the required image exists, and fetch the content of the image.
type Local interface {
	ImageExist(imageName string) (bool, error)
	ImageSave(imageName string) (io.ReadCloser, error)
}

// Interface Remote means that the container runtime is running at the "server side"/"kubelet side".
// Interface Remote load the content of the image into the container runtime.
type Remote interface {
	LoadImage(content io.ReadCloser) error
}
