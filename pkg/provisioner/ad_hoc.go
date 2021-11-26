package provisioner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

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

const waitPodRunningInterval = 2 * time.Second
const waitPodRunningTimeout = 5 * time.Minute
const imageKubectlPushPeer = "ghcr.io/strrl/kubectl-push-peer:latest"
const defaultKubectlPushPort = 28375

// TODO: provisioner interface

// Struct AdHoc will set up a temporary Pod (with image: kubectl-push-peer) on the specified node.
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

func (it *AdHoc) SpawnPeerOnTargetNode(ctx context.Context, node string) (Peer, error) {
	podName := fmt.Sprintf("kubectl-push-peer-on-%s", node)

	// if the pod already exists, delete it
	if _, err := it.clientset.CoreV1().Pods(it.namespace).Get(ctx, podName, metav1.GetOptions{}); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, err
		}
	} else {
		getLogger().WithName("ad-hoc").Info("Pod already existed, delete it", "pod", podName)
		if err := it.clientset.CoreV1().Pods(it.namespace).Delete(ctx, podName, metav1.DeleteOptions{}); err != nil {
			return nil, err
		}
		// wait for the pod to be deleted
		waitDeleteErr := wait.PollImmediate(waitPodRunningInterval, waitPodRunningTimeout, func() (bool, error) {
			_, err := it.clientset.CoreV1().Pods(it.namespace).Get(ctx, podName, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					return true, nil
				}
				return false, err
			}
			return false, nil
		})
		if waitDeleteErr != nil {
			return nil, waitDeleteErr
		}
	}

	getLogger().WithName("ad-hoc").Info("create Pod with kubectl-push-peer", "pod", podName)

	_, err := it.clientset.CoreV1().Pods(it.namespace).Create(ctx, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: it.namespace,
		},
		Spec: v1.PodSpec{
			NodeName: node,
			Containers: []v1.Container{
				{
					Name:            "kubectl-push-peer",
					Image:           imageKubectlPushPeer,
					ImagePullPolicy: v1.PullAlways,
					ReadinessProbe: &v1.Probe{
						Handler: v1.Handler{
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
		return nil, err
	}

	var pod *v1.Pod

	getLogger().WithName("ad-hoc").Info("waiting for pod start up", "pod", podName)
	// wait for pod to be running
	err = wait.Poll(waitPodRunningInterval, waitPodRunningTimeout, func() (done bool, err error) {
		pod, err = it.clientset.CoreV1().Pods(it.namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return podReady(&pod.Status), nil
	})

	if err != nil {
		return nil, err
	}

	// port forward, restore the context
	getLogger().WithName("ad-hoc").Info("setup port-forwarding", "pod", podName)
	portForwardCtx, cancelFunc := context.WithCancel(ctx)
	localPort, err := it.portForward(portForwardCtx, pod, it.restconfig, defaultKubectlPushPort)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	// return AdHocPeer
	result := &adHocPeer{
		nodeName:              node,
		pod:                   pod,
		portForwardCancelFunc: cancelFunc,
		clientset:             it.clientset,
		localPort:             localPort,
	}
	return result, nil
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
	// pod is the corresponding kubectl-push-peer pod.
	pod       *v1.Pod
	clientset *kubernetes.Clientset
	// localPort is the listen port on the local machine after port forwarding
	localPort uint16
}

// Destroy would terminate the port forwarding, and delete the ad-hoc kubectl-push-peer pod.
func (it *adHocPeer) Destory() error {
	it.portForwardCancelFunc()
	return it.clientset.CoreV1().Pods(it.pod.Namespace).Delete(context.TODO(), it.pod.Name, metav1.DeleteOptions{})
}

func (it *adHocPeer) BaseUrl() string {
	return fmt.Sprintf("http://localhost:%d", it.localPort)
}

func (it *AdHoc) portForward(ctx context.Context, pod *v1.Pod, restconfig *rest.Config, remotePort uint16) (uint16, error) {
	// kubernetes port forward
	req := it.clientset.CoreV1().RESTClient().Post().Resource("pods").Namespace(pod.Namespace).Name(pod.Name).SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(restconfig)
	if err != nil {
		return 0, err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())

	preader, pwriter := io.Pipe()

	go func() {
		<-ctx.Done()
		pwriter.Close()
	}()
	go func() {
		// TODO: forward the logs from port forwarder
		io.Copy(io.Discard, preader)
	}()

	readyChan := make(chan struct{})
	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("0:%d", remotePort)}, ctx.Done(), readyChan, pwriter, pwriter)
	if err != nil {
		return 0, err
	}

	errChan := make(chan error)
	go func() {
		errChan <- forwarder.ForwardPorts()
	}()
	select {
	case <-readyChan:
		forwardedPorts, err := forwarder.GetPorts()
		if err != nil {
			return 0, nil
		}

		// return the first forwarded port
		return forwardedPorts[0].Local, nil
	case err := <-errChan:
		return 0, err
	}
}
