package provisioner

// Agent means one "kubectl-image-agent" instance that runs on specific node.
// It should be cleaned up after it is no longer needed.
// You could build client.HTTPClient with Agent.BaseURL().
type Agent interface {
	// Destroy the agent, release the resources.
	// It MUST be called when the agent is no longer needed.
	Destroy() error

	// TODO: change the name. BaseURL sounds not so good. :(
	BaseURL() string
}
