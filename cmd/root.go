package cmd

import (
	"os"

	version_command "github.com/diogopereiradev/httpzen/cmd/commands/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "httpzen",
	Short: "HTTP Zen CLI Tool for API Management and Development",
}

func init() {
	// Commands
	version_command.Executor(rootCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
