package cmd

import (
	"github.com/STRRL/kubectl-image/pkg/version"
	"github.com/spf13/cobra"
)

// NewRootCommand is the constructor for root command.
func NewRootCommand() *cobra.Command {
	result := &cobra.Command{
		Use:     "kubectl-image",
		Example: "kubectl-image -h",
		Short:   "docker image but for kubernetes",
		Version: version.GetVersion(),
	}
	result.AddCommand(
		NewLoadCommand(),
		NewListCommand(),
	)

	return result
}
