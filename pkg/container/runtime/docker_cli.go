package runtime

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

var _ Remote = (*DockerCli)(nil)
var _ Local = (*DockerCli)(nil)

// DockerCli introduce image operations with docker.
// Deprecated: use Docker instead.
type DockerCli struct{}

// LoadImage implements Remote.LoadImage.
func (it *DockerCli) LoadImage(content io.ReadCloser) error {
	command := exec.Command("docker", "image", "load")
	command.Stdin = content
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Start(); err != nil {
		return errors.Wrap(err, "execute docker image load")
	}

	if err := command.Wait(); err != nil {
		return errors.Wrap(err, "wait docker image load")
	}

	if command.ProcessState.ExitCode() != 0 {
		return errors.Errorf("exit code is not 0, exitcode: %d", command.ProcessState.ExitCode())
	}

	return nil
}

// ImageExist implements Local.ImageExist.
func (it *DockerCli) ImageExist(imageName string) (bool, error) {
	command := exec.Command("docker", "image", "inspect", imageName)
	if err := command.Start(); err != nil {
		return false, errors.Wrap(err, "execute docker image inspect")
	}

	if err := command.Wait(); err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("no such image %s", imageName))
	}

	return true, nil
}

// ImageSave implements Local.ImageSave.
func (it *DockerCli) ImageSave(imageName string, content io.Writer) error {
	command := exec.Command("docker", "image", "save", imageName)
	command.Stdout = content

	if err := command.Start(); err != nil {
		return errors.Wrap(err, "execute docker image save command")
	}

	if err := command.Wait(); err != nil {
		return errors.Wrap(err, "wait docker image save")
	}

	if command.ProcessState.ExitCode() != 0 {
		return errors.Errorf("exit code is not 0, exitcode: %d", command.ProcessState.ExitCode())
	}

	return nil
}
