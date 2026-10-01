package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/oas-cli/internal/formatter"
	"github.com/yourusername/oas-cli/internal/service"
)

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func inspectHandler(filename string) error {
	result, err := service.Inspect(filename)
	if err != nil {
		return err
	}
	formatter.FormatInspect(os.Stdout, result)
	return nil
}

// definition for the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Display an overview of an OpenAPI specification",
	Long:  `Inspect an OpenAPI specification file and display key information about the API, including its version, metadata, available endpoints, and defined schemas.`,
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		return inspectHandler(args[0])
	},
}
