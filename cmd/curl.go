package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/oas-cli/internal/formatter"
	"github.com/yourusername/oas-cli/internal/service"
)

func init() {
	rootCmd.AddCommand(curlCmd)
}

func curlHandler(filename string, method string, path string) error {
	result, err := service.Curl(filename, method, path)
	if err != nil {
		return err
	}
	formatter.FormatCurl(os.Stdout, result)
	return nil
}

var curlCmd = &cobra.Command{
	Use:   "curl",
	Short: "Generates an example curl request for an operation",
	Long:  `Generates an example curl command for one operation in the OpenAPI spec, including the server URL, path/query parameters, and an example request body if defined.`,
	Args:  cobra.ExactArgs(3),

	RunE: func(cmd *cobra.Command, args []string) error {
		return curlHandler(args[0], args[1], args[2])
	},
}
