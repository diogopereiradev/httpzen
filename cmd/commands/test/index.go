package test_command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Executor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command for testing the application functionality",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Executed test command")
		},
	}
	return cmd
}
