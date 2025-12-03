package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the CLI
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hf-go",
		Short: "Hugging Face Models CLI",
		Long:  `A CLI tool for interacting with Hugging Face models.`,
	}

	// Add subcommands
	cmd.AddCommand(NewListModelsCmd())

	return cmd
}

// Execute runs the CLI
func Execute() error {
	return NewRootCmd().Execute()
}
