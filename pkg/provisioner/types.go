package provisioner

// Interace Peer means one "kubectl-push-peer" instance that runs on specific node.
// It should be cleaned up after it is no longer needed.
type Peer interface {
	Destory() error
	BaseUrl() string
}
