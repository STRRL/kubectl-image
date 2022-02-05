package cmd

import (
	"context"
	"io"
	"os"

	"github.com/STRRL/kubectl-image/pkg/agent"
	"github.com/STRRL/kubectl-image/pkg/agent/provisioner"
	"github.com/STRRL/kubectl-image/pkg/util"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewLoadCommand is the constructor for command load.
func NewLoadCommand(logger logr.Logger) *cobra.Command {
	loadOption := LoadCommandOptions{}

	result := &cobra.Command{
		Use:   "load",
		Short: "Load an image from a tar archive or STDIN",
		Run: func(cmd *cobra.Command, args []string) {
			var inputStream io.ReadCloser
			defer func() {
				if inputStream != nil {
					err := inputStream.Close()
					if err != nil {
						logger.Error(err, "close input stream")
					}
				}
			}()

			if len(loadOption.Input) == 0 {
				inputStream = os.Stdin
			} else {
				file, err := os.Open(loadOption.Input)
				if err != nil {
					logger.Error(err, "open input file", "filename", loadOption.Input)

					return
				}
				inputStream = file
			}
			clientsForNodes, err := spawnClientsForEachNode(context.TODO(), logger)
			if err != nil {
				logger.Error(err, "spawn clients for each node")

				return
			}
			for _, client := range clientsForNodes {
				err := client.LoadImage(context.TODO(), inputStream)
				if err != nil {
					logger.Error(err, "load image")

					return
				}
			}
		},
	}

	result.Flags().StringVarP(
		&(loadOption.Input),
		"input",
		"i",
		"",
		"Read from tar archive file, instead of STDIN")

	return result
}

// LoadCommandOptions is the options/flags for command load.
type LoadCommandOptions struct {
	Input string
}

func spawnClientsForEachNode(ctx context.Context, logger logr.Logger) ([]agent.Client, error) {
	var result []agent.Client

	clientset, restConfig, rawConfig, err := util.LoadClientsetAndConfiguration()
	if err != nil {
		return nil, errors.Wrap(err, "load clientset and configuration")
	}

	adHocProvisioner := provisioner.NewAdHoc(rawConfig.Contexts[rawConfig.CurrentContext].Namespace, clientset, restConfig)

	nodeList, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "list nodes")
	}

	for _, node := range nodeList.Items {
		spawnedAgent, err := adHocProvisioner.SpawnPeerOnTargetNode(ctx, node.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "spawn agent on node %s", node.Name)
		}

		client := agent.NewHTTPClient(spawnedAgent.BaseURL(), logger)
		result = append(result, client)
	}

	return result, nil
}
