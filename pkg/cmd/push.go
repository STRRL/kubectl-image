package cmd

import (
	"context"
	"io"

	containerruntime "github.com/STRRL/kubectl-push/pkg/container/runtime"
	"github.com/STRRL/kubectl-push/pkg/peer"
	"github.com/STRRL/kubectl-push/pkg/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type CmdPushOptions struct {
	configFlags *genericclioptions.ConfigFlags
	image       string
}

// constructor for CmdPushOptions
func NewCmdPushOptions() *CmdPushOptions {
	return &CmdPushOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

// Run executes the command
func (o *CmdPushOptions) RunE() error {
	cr := containerruntime.Docker{}
	var err error
	var exist bool
	if exist, err = cr.ImageExist(o.image); err != nil {
		return err
	}

	if !exist {
		return errors.Errorf("Image %s does not exist on local machine", o.image)
	}

	// prepare kubectl-push-peer
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, nil)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return err
	}

	peerProvisioner := provisioner.NewAdHoc(rawConfig.Contexts[rawConfig.CurrentContext].Namespace, clientset, restConfig)
	ctx := context.TODO()

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, node := range nodes.Items {
		preader, pwriter := io.Pipe()
		go func() {
			// TODO: handle these errors
			if err := cr.ImageSave(o.image, pwriter); err != nil {
				getLogger().Error(err, "failed to save image", "image", o.image)
			}
			pwriter.Close()
			getLogger().Info("image saved", "image", o.image, "node", node.Name)
		}()

		peerInstance, err := peerProvisioner.SpawnPeerOnTargetNode(ctx, node.Name)
		if err != nil {
			return err
		}
		defer peerInstance.Destory()

		getLogger().Info("image transmitting", "image", o.image, "node", node.Name)
		baseUrl := peerInstance.BaseUrl()
		if err := peer.LoadImage(ctx, baseUrl, preader); err != nil {
			return err
		}
	}
	return nil
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
