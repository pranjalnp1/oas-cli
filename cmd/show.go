package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/oas-cli/internal/service"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

func showHandler(filename string, method string, path string, writer io.Writer) error {
	return service.Show(filename, method, path, writer)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays detailed information about a single operation",
	Long:  `Displays detailed information about one operation in the OpenAPI spec, including its summary, parameters, and responses.`,
	Args:  cobra.ExactArgs(3),

	RunE: func(cmd *cobra.Command, args []string) error {
		return showHandler(args[0], args[1], args[2], os.Stdout)
	},
}
