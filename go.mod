module github.com/STRRL/kubectl-push

go 1.16

require (
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zapr v1.2.2
	github.com/google/uuid v1.2.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.20.0
	k8s.io/api v0.23.1
	k8s.io/apimachinery v0.23.1
	k8s.io/cli-runtime v0.23.1
	k8s.io/client-go v0.23.1
)

replace k8s.io/klog/v2 => k8s.io/klog/v2 v2.20.0
