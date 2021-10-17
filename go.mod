module github.com/STRRL/kubectl-push

go 1.16

require (
	github.com/go-logr/zapr v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	go.uber.org/zap v1.19.0
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/cli-runtime v0.22.2
	k8s.io/client-go v0.22.2
)

replace k8s.io/klog/v2 => k8s.io/klog/v2 v2.20.0
