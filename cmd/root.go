package cmd

import (
	"os"

	help_command "github.com/diogopereiradev/httpzen/cmd/commands/help"
	version_command "github.com/diogopereiradev/httpzen/cmd/commands/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "httpzen",
	Short: "Httpzen CLI Tool for API Management and Development",
}

func init() {
	help_command.Init(rootCmd)
	version_command.Init(rootCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
