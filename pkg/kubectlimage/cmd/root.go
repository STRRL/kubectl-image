package cmd

import "github.com/spf13/cobra"

// NewRootCommand is the constructor for root command.
func NewRootCommand() *cobra.Command {
	result := &cobra.Command{
		Use:     "kubectl-image",
		Example: "kubectl-image -h",
		Short:   "docker image but for kubernetes",
		Version: "",
	}
	result.AddCommand(
		NewLoadCommand(),
		NewListCommand(),
	)

	return result
}
