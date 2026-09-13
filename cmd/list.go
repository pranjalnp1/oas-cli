package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/oas-cli/internal/service"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func listHandler(filename string, writer io.Writer) error { return service.List(filename, writer) }

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists every operation (method + path) defined in the openApi spec",
	Long:  `Lists every operation (method + path) defined in the openApi spec.`,
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error { return listHandler(args[0], os.Stdout) },
}
