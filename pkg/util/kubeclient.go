package util

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// LoadClientsetAndConfiguration would load configuration with kubectl convention.
func LoadClientsetAndConfiguration() (
	*kubernetes.Clientset,
	*rest.Config,
	*api.Config,
	error,
) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "load rest config")
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "setup kubeClient config")
	}

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "fetch rawConfig from clientConfig")
	}

	return clientset, restConfig, &rawConfig, nil
}
