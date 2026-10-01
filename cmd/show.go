package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/oas-cli/internal/formatter"
	"github.com/yourusername/oas-cli/internal/service"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

func showHandler(filename string, method string, path string) error {
	result, err := service.Show(filename, method, path)
	if err != nil {
		return err
	}
	formatter.FormatShow(os.Stdout, result)
	return nil
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays detailed information about a single operation",
	Long:  `Displays detailed information about one operation in the OpenAPI spec, including its summary, parameters, and responses.`,
	Args:  cobra.ExactArgs(3),

	RunE: func(cmd *cobra.Command, args []string) error {
		return showHandler(args[0], args[1], args[2])
	},
}
