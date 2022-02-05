package provisioner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	waitPodRunningInterval = 2 * time.Second
	waitPodRunningTimeout  = 5 * time.Minute
	imageKubectlPushPeer   = "ghcr.io/strrl/kubectl-image-agent:latest"
	defaultKubectlPushPort = 28375
)

// TODO: provisioner interface

// AdHoc will set up a temporary Pod (with image: kubectl-image-agent) on the specified node.
type AdHoc struct {
	namespace  string
	clientset  *kubernetes.Clientset
	restconfig *rest.Config
}

// NewAdHoc creates a new AdHoc instance.
func NewAdHoc(namespace string, clientset *kubernetes.Clientset, restconfig *rest.Config) *AdHoc {
	return &AdHoc{
		namespace:  namespace,
		clientset:  clientset,
		restconfig: restconfig,
	}
}

// SpawnPeerOnTargetNode would initialize agent on all the kubernetes nodes.
// TODO: the interface of provisioner.
func (it *AdHoc) SpawnPeerOnTargetNode(ctx context.Context, node string) (Peer, error) {
	podName := fmt.Sprintf("kubectl-image-agent-on-%s", node)

	if err := it.deletePeerIfAlreadyExists(ctx, podName); err != nil {
		return nil, errors.Wrap(err, "delete agent if already exists")
	}

	if err := it.spawnNewPeerOnTargetNode(ctx, node, podName); err != nil {
		return nil, errors.Wrapf(err, "spawn new agent on target node %s, podName %s", node, podName)
	}

	if err := it.waitNewPeerIsReady(ctx, podName); err != nil {
		return nil, errors.Wrapf(err, "wait new agent pod %s is ready", podName)
	}

	localPort, cancelFunc, err := it.establishPortForward(ctx, podName)
	if err != nil {
		return nil, errors.Wrap(err, "establish port forward")
	}

	pod, err := it.clientset.CoreV1().Pods(it.namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "fetch pod %s", podName)
	}

	return &adHocPeer{
		nodeName:              node,
		pod:                   pod,
		portForwardCancelFunc: cancelFunc,
		clientset:             it.clientset,
		localPort:             localPort,
	}, nil
}

func (it *AdHoc) deletePeerIfAlreadyExists(ctx context.Context, podName string) error {
	// if the pod already exists, delete it
	_, err := it.clientset.CoreV1().Pods(it.namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}

		return errors.Wrapf(err, "fetch existed kubectl-image-agent pod")
	}

	getLogger().WithName("ad-hoc").Info("Pod already existed, delete it", "pod", podName)

	err = it.clientset.
		CoreV1().
		Pods(it.namespace).
		Delete(ctx, podName, metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "delete pod %s/%s", it.namespace, podName)
	}
	// wait for the pod to be deleted
	waitDeleteErr := wait.PollImmediate(waitPodRunningInterval, waitPodRunningTimeout, func() (bool, error) {
		_, err := it.clientset.CoreV1().Pods(it.namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}

			return false, errors.Wrapf(err, "fetch status of pod %s", podName)
		}

		return false, nil
	})

	if waitDeleteErr != nil {
		return errors.Wrapf(waitDeleteErr, "wait for pod %s to be deleted", podName)
	}

	return nil
}

func (it *AdHoc) spawnNewPeerOnTargetNode(ctx context.Context, node, podName string) error {
	getLogger().WithName("ad-hoc").Info("create Pod with kubectl-image-agent", "pod", podName)

	_, err := it.clientset.CoreV1().Pods(it.namespace).Create(ctx, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: it.namespace,
		},
		Spec: v1.PodSpec{
			NodeName: node,
			Containers: []v1.Container{
				{
					Name:            "kubectl-image-agent",
					Image:           imageKubectlPushPeer,
					ImagePullPolicy: v1.PullAlways,
					ReadinessProbe: &v1.Probe{
						ProbeHandler: v1.ProbeHandler{
							HTTPGet: &v1.HTTPGetAction{
								Path: "/healthz",
								Port: intstr.FromInt(defaultKubectlPushPort),
							},
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "create kube-push-agent pod on node %s", node)
	}

	return nil
}

func (it *AdHoc) waitNewPeerIsReady(ctx context.Context, podName string) error {
	var pod *v1.Pod

	getLogger().WithName("ad-hoc").Info("waiting for pod start up", "pod", podName)
	// wait for pod to be running
	err := wait.Poll(waitPodRunningInterval, waitPodRunningTimeout, func() (done bool, err error) {
		pod, err = it.clientset.CoreV1().Pods(it.namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return false, errors.Wrapf(err, "fetch pod %s", podName)
		}

		return podReady(&pod.Status), nil
	})
	if err != nil {
		return errors.Wrapf(err, "wait for pod %s to be ready", podName)
	}

	return nil
}

func (it *AdHoc) establishPortForward(ctx context.Context, podName string) (uint16, context.CancelFunc, error) {
	pod, err := it.clientset.CoreV1().Pods(it.namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return 0, nil, errors.Wrapf(err, "port forard for pod %s", podName)
	}

	getLogger().WithName("ad-hoc").Info("setup port-forwarding", "pod", podName)

	portForwardCtx, cancelFunc := context.WithCancel(ctx)

	localPort, err := it.portForward(portForwardCtx, pod, it.restconfig, defaultKubectlPushPort)
	if err != nil {
		cancelFunc()

		return 0, nil, err
	}

	return localPort, cancelFunc, nil
}

func podReady(podStatus *v1.PodStatus) bool {
	if podStatus == nil {
		return false
	}

	for _, condition := range podStatus.Conditions {
		if condition.Type == v1.PodReady {
			return condition.Status == v1.ConditionTrue
		}
	}

	return false
}

type adHocPeer struct {
	nodeName string
	// portForwardCancelFunc is used to cancel the port forwarding.
	portForwardCancelFunc context.CancelFunc
	// pod is the corresponding kubectl-image-agent pod.
	pod       *v1.Pod
	clientset *kubernetes.Clientset
	// localPort is the listen port on the local machine after port forwarding
	localPort uint16
}

// Destroy would terminate the port forwarding, and delete the ad-hoc kubectl-image-agent pod.
func (it *adHocPeer) Destroy() error {
	it.portForwardCancelFunc()
	err := it.clientset.CoreV1().Pods(it.pod.Namespace).Delete(context.TODO(), it.pod.Name, metav1.DeleteOptions{})

	return errors.Wrapf(err, "delete pod %s/%s", it.pod.Namespace, it.pod.Name)
}

func (it *adHocPeer) BaseURL() string {
	return fmt.Sprintf("http://localhost:%d", it.localPort)
}

func (it *AdHoc) portForward(
	ctx context.Context, pod *v1.Pod, restconfig *rest.Config, remotePort uint16,
) (uint16, error) {
	// kubernetes port forward
	req := it.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(pod.Namespace).Name(pod.Name).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(restconfig)
	if err != nil {
		return 0, errors.Wrapf(err, "build troundtrip transport and upgrader")
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())
	preader, pwriter := io.Pipe()

	go func() {
		<-ctx.Done()

		_ = pwriter.Close()
	}()
	go func() {
		// TODO: forward the logs from port forwarder
		if _, err := io.Copy(os.Stderr, preader); err != nil {
			getLogger().Error(err, "forward logs from port forwarder")
		}
	}()

	readyChan := make(chan struct{})

	forwarder, err := portforward.New(
		dialer,
		[]string{
			fmt.Sprintf("0:%d", remotePort),
		},
		ctx.Done(), readyChan, pwriter, pwriter,
	)
	if err != nil {
		return 0, errors.Wrapf(err, "build port forwarder")
	}

	errChan := make(chan error)

	go func() {
		errChan <- forwarder.ForwardPorts()
	}()
	select {
	case <-readyChan:
		forwardedPorts, err := forwarder.GetPorts()
		if err != nil {
			return 0, errors.Wrapf(err, "get forwarded ports")
		}

		// return the first forwarded port
		return forwardedPorts[0].Local, nil
	case err := <-errChan:
		return 0, errors.Wrapf(err, "error from forwarder.ForwardPorts")
	}
}
