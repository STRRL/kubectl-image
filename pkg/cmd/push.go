package cmd

import (
	"context"
	"io"

	containerruntime "github.com/STRRL/kubectl-push/pkg/container/runtime"
	"github.com/STRRL/kubectl-push/pkg/peer"
	"github.com/STRRL/kubectl-push/pkg/provisioner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	var cr containerruntime.Local
	var err error
	var exist bool
	if exist, err = cr.ImageExist(o.image); err != nil {
		return err
	}

	if !exist {
		return errors.Errorf("Image %s does not exist on local machine", o.image)
	}

	var reader io.ReadCloser
	if reader, err = cr.ImageSave(o.image); err != nil {
		return err
	}
	defer reader.Close()

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

	peerProvisioner := provisioner.NewAdHoc(rawConfig.Contexts[rawConfig.CurrentContext].Namespace, clientset)
	ctx := context.TODO()

	peerInstance, err := peerProvisioner.SpawnPeerOnTargetNode(ctx, "nodeName")
	if err != nil {
		return err
	}
	defer peerInstance.Destory()

	baseUrl := peerInstance.BaseUrl()
	if err := peer.LoadImage(ctx, baseUrl, reader); err != nil {
		return nil
	}
	return nil
}

func NewCmdPush() *cobra.Command {
	o := NewCmdPushOptions()

	cmd := &cobra.Command{
		Use:          "push",
		Short:        "Push an image to kubernetes nodes",
		Example:      "push alpine:latest",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			options := NewCmdPushOptions()
			return options.RunE()
		},
	}

	cmd.Flags().StringVarP(&o.image, "image", "i", "", "Image to push")

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}
