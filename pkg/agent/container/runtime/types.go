package runtime

import "io"

// Local means that the container runtime is running at the "client side"/"kubectl side".
// Local check the required image exists, and fetch the content of the image.
type Local interface {
	ImageExist(imageName string) (bool, error)
	ImageSave(imageName string, content io.Writer) error
}

// Remote means that the container runtime is running at the "server side"/"kubelet side".
// Remote load the content of the image into the container runtime.
type Remote interface {
	// LoadImage loads image from bytes.
	LoadImage(content io.ReadCloser) error
}
