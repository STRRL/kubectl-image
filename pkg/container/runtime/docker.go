package runtime

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

var _ Remote = (*Docker)(nil)

type Docker struct {
}

func (it *Docker) LoadImage(content io.ReadCloser) error {
	command := exec.Command("docker", "image", "load")
	command.Stdin = content
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		return err
	}
	if err := command.Wait(); err != nil {
		return err
	}
	if command.ProcessState.ExitCode() != 0 {
		return errors.Errorf("exit code is not 0, exitcode: %d", command.ProcessState.ExitCode())
	}
	return nil
}

func (it *Docker) ImageExist(imageName string) (bool, error) {
	// TODO: implement
	return true, nil
}

func (it *Docker) ImageSave(imageName string, content io.Writer) error {
	command := exec.Command("docker", "image", "save", imageName)
	command.Stdout = content

	if err := command.Start(); err != nil {
		return err
	}
	if err := command.Wait(); err != nil {
		return err
	}
	if command.ProcessState.ExitCode() != 0 {
		return errors.Errorf("exit code is not 0, exitcode: %d", command.ProcessState.ExitCode())
	}
	return nil
}
