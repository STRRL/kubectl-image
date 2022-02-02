package provisioner

// Peer means one "kubectl-push-peer" instance that runs on specific node.
// It should be cleaned up after it is no longer needed.
type Peer interface {
	// Destroy the peer, release the resources.
	// It MUST be called when the peer is no longer needed.
	Destroy() error

	// TODO: change the name. BaseURL sounds not so good. :(
	BaseURL() string
}
