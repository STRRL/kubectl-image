package cmd

import (
	"context"
	"fmt"
	"io"

	containerruntime "github.com/STRRL/kubectl-push/pkg/container/runtime"
	"github.com/STRRL/kubectl-push/pkg/peer"
	"github.com/STRRL/kubectl-push/pkg/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type PushCommandOptions struct {
	configFlags *genericclioptions.ConfigFlags
	image       string
}

// NewCmdPushOptions is the constructor for PushCommandOptions.
func NewCmdPushOptions() *PushCommandOptions {
	return &PushCommandOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

// RunE executes the command.
func (o *PushCommandOptions) RunE() error {
	containerRuntime := containerruntime.Docker{}

	var (
		err   error
		exist bool
	)

	if exist, err = containerRuntime.ImageExist(o.image); err != nil {
		return errors.Wrap(err, "check image exists")
	}

	if !exist {
		return errors.Errorf("Image %s does not exist on local machine", o.image)
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return errors.Wrap(err, "load rest config")
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return errors.Wrap(err, "setup kubeClient config")
	}

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return errors.Wrap(err, "fetch rawConfig from clientConfig")
	}

	peerProvisioner := provisioner.NewAdHoc(rawConfig.Contexts[rawConfig.CurrentContext].Namespace, clientset, restConfig)
	ctx := context.TODO()

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "list nodes")
	}

	peers, err := o.preparePeersOnEachNode(ctx, nodes, containerRuntime, err, peerProvisioner)
	if err != nil {
		return err
	}

	defer func() {
		for _, item := range peers {
			peerInstance := item.peer
			nodeName := item.node

			if err := peerInstance.Destroy(); err != nil {
				getLogger().Error(err, "destroy peer instance", "node", nodeName)
			}
		}
	}()

	return nil
}

type peerAndNodeName struct {
	peer provisioner.Peer
	node string
}

func (o *PushCommandOptions) preparePeersOnEachNode(
	ctx context.Context,
	nodes *v1.NodeList,
	containerRuntime containerruntime.Docker,
	err error,
	peerProvisioner *provisioner.AdHoc,
) ([]peerAndNodeName, error) {
	var peers []peerAndNodeName

	for _, node := range nodes.Items {
		preader, pwriter := io.Pipe()

		go func() {
			// TODO: handle these errors
			if err := containerRuntime.ImageSave(o.image, pwriter); err != nil {
				getLogger().Error(err, "failed to save image", "image", o.image)
			}

			err = pwriter.Close()
			if err != nil {
				getLogger().Error(err, "close pipe writer")
			}

			getLogger().Info("image saved", "image", o.image, "node", node.Name)
		}()

		peerInstance, err := peerProvisioner.SpawnPeerOnTargetNode(ctx, node.Name)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("spawn peer on node %s", node.Name))
		}

		peers = append(peers, peerAndNodeName{
			peer: peerInstance,
			node: node.Name,
		})

		getLogger().Info("image transmitting", "image", o.image, "node", node.Name)

		baseURL := peerInstance.BaseURL()
		if err := peer.LoadImage(ctx, baseURL, preader); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("load image for node %s", node.Name))
		}
	}

	return peers, nil
}

func NewCmdPush() *cobra.Command {
	options := NewCmdPushOptions()

	cmd := &cobra.Command{
		Use:          "push",
		Short:        "Push an image to kubernetes nodes",
		Example:      "push alpine:latest",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.RunE()
		},
	}

	cmd.Flags().StringVarP(&options.image, "image", "i", "", "Image to push")

	options.configFlags.AddFlags(cmd.Flags())

	return cmd
}
