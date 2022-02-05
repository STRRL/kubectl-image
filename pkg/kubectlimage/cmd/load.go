package cmd

import "github.com/spf13/cobra"

// NewLoadCommand is the constructor for command load.
func NewLoadCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "load",
		Short: "Load an image from a tar archive or STDIN",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return result
}
