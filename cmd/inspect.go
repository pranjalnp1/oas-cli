package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/oas-cli/internal/service"
)

func init() {
	rootCmd.AddCommand(inspectCmd)
}

// definition for the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Display an overview of an OpenAPI specification",
	Long:  `Inspect an OpenAPI specification file and display key information about the API, including its version, metadata, available endpoints, and defined schemas.`,
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		return service.Inspect(args[0])
	},
}
