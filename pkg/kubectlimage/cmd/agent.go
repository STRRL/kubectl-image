package cmd

import "github.com/spf13/cobra"

// NewAgentCommand is the constructor for command agent.
func NewAgentCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "agent",
		Short: "Manage the kubectl-image agent",
	}
	result.AddCommand(
		NewAgentPrepareCommand(),
		NewAgentCleanupCommand(),
	)

	return result
}

// NewAgentPrepareCommand is the constructor for command agent prepare.
func NewAgentPrepareCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare the kubectl-image agent",
	}

	return result
}

// NewAgentCleanupCommand is the constructor for command agent cleanup.
func NewAgentCleanupCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup the kubectl-image agent",
	}

	return result
}
