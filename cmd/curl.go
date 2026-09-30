package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/oas-cli/internal/formatter"
	"github.com/yourusername/oas-cli/internal/service"
)

var curlBaseURL string

func init() {
	curlCmd.Flags().StringVar(&curlBaseURL, "base-url", "", "override the server URL (required if the spec's server URL is relative)")
	rootCmd.AddCommand(curlCmd)
}

func curlHandler(filename string, method string, path string, baseURL string) error {
	result, err := service.Curl(filename, method, path, baseURL)
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
		return curlHandler(args[0], args[1], args[2], curlBaseURL)
	},
}
