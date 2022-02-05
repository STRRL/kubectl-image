package cmd

import "github.com/spf13/cobra"

// NewListCommand is the constructor for command list.
func NewListCommand() *cobra.Command {
	result := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List images",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return result
}
