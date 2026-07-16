package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
		fmt.Println("inspecting openApi spec...")
		return nil
	},
}
