package cmd

import (
	"os"

	test_command "github.com/diogopereiradev/httpzen/cmd/commands/test"
	version_flag "github.com/diogopereiradev/httpzen/cmd/flags/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "httpzen",
	Short: "HTTP Zen CLI Tool for API Management and Development",
}

func init() {
	// Commands
	rootCmd.AddCommand(test_command.Executor())

	// Flags
	version_flag.AddFlag(rootCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
