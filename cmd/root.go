package cmd

import (
	"os"

	clean_cache_command "github.com/diogopereiradev/httpzen/cmd/commands/clean-cache"
	help_command "github.com/diogopereiradev/httpzen/cmd/commands/help"
	request_command "github.com/diogopereiradev/httpzen/cmd/commands/request"
	version_command "github.com/diogopereiradev/httpzen/cmd/commands/version"
	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "httpzen [METHOD] [URL] [PARAMETERS...]",
	Short: "Httpzen CLI Tool for API Management and Development",
	Args:  cobra.ArbitraryArgs,
}

func setFlagErrorFunc(cmd *cobra.Command) {
	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		rootCmd.Help()
		if err != nil {
			logger_module.Error(err.Error(), 50)
		}
		return nil
	})
	for _, sub := range cmd.Commands() {
		setFlagErrorFunc(sub)
	}
}

func init() {
	help_command.Init(rootCmd)
	version_command.Init(rootCmd)
	request_command.Init(rootCmd)
	clean_cache_command.Init(rootCmd)

	setFlagErrorFunc(rootCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
