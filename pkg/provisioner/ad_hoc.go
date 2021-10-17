package provisioner

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const imageKubectlPushPeer = "ghcr.io/strrl/kubectl-push-peer:latest"

// TODO: provisioner interface

// Struct AdHoc will set up a temporary Pod (with image: kubectl-push-peer) on the specified node.
type AdHoc struct {
	namespace string
	clientset *kubernetes.Clientset
}

// NewAdHoc creates a new AdHoc instance.
func NewAdHoc(namespace string, clientset *kubernetes.Clientset) *AdHoc {
	return &AdHoc{
		namespace: namespace,
		clientset: clientset,
	}
}

func (it *AdHoc) SpawnPeerOnTargetNode(ctx context.Context, node string) (Peer, error) {
	_, err := it.clientset.CoreV1().Pods(it.namespace).Create(ctx, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubectl-push-peer",
			Namespace: it.namespace,
		},
		Spec: v1.PodSpec{
			NodeName: node,
			Containers: []v1.Container{
				{
					Name:  "kubectl-push-peer",
					Image: imageKubectlPushPeer,
				},
			},
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	// wait for pod to be running

	// port forward, restore the context

	// return AdHocPeer

	panic("wip")
}

func (it *AdHoc) destoryPeer(peer Peer) error {
	// cancel port forward

	// delete pod

	panic("wip")
}

type adHocPeer struct {
	nodeName     string
	controlledBy *AdHoc
}

func (it *adHocPeer) Destory() error {
	return it.controlledBy.destoryPeer(it)
}

func (it *adHocPeer) BaseUrl() string {
	panic("wip")
}

func portForward(ctx context.Context, pod *v1.Pod, localPort int, remotePort int) error {

	// kubernetes port forward

	panic("wip")
}
