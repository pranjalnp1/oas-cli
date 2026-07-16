package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// definition for the root command
var rootCmd = &cobra.Command{
	Use:   `oas`,
	Short: `A CLI for inspecting OpenAPI specifications`,
	Long: `oas is a command-line tool for inspecting OpenAPI 3
		specifications. It can list endpoints, inspect operations,
		and generate example requests.`,
}

// Execute function for main.go to use as an entry point
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
