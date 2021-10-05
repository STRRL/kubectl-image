package importer

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type DockerImageLoader struct {
}

func (it *DockerImageLoader) LoadImage(content io.ReadCloser) error {
	command := exec.Command("docker", "image", "load")
	pipe, err := command.StdinPipe()
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err != nil {
		return err
	}
	err = command.Start()
	if err != nil {
		return err
	}
	copied, err := io.Copy(pipe, content)
	logger.Info("image content copied", "copied-bytes", copied)
	if err != nil {
		return err
	}
	pipe.Close()
	err = command.Wait()
	if err != nil {
		return err
	}
	if command.ProcessState.ExitCode() != 0 {
		return errors.Errorf("exit code is not 0, exitcode: %d", command.ProcessState.ExitCode())
	}
	return nil
}
