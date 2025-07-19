package test_command

import (
	"fmt"

	app_config "github.com/diogopereiradev/httpzen/internal/config"
	"github.com/spf13/cobra"
)

func Executor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command for testing the application functionality",
		Run: func(cmd *cobra.Command, args []string) {
			config := app_config.GetConfig()
			fmt.Println("Current Slow Response Threshold:", config.SlowResponseThreshold, "ms")
			fmt.Println("Executed test command")
		},
	}
	return cmd
}
