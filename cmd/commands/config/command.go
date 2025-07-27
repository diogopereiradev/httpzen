package config_command

import (
	"fmt"

	"github.com/diogopereiradev/httpzen/internal/components/config_editor"
	"github.com/spf13/cobra"
)

var runConfigEditor = config_editor.RunConfigEditor

func Init(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage the app configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runConfigEditor(); err != nil {
				fmt.Println("Error on execution of editor component:", err)
			}
		},
	}
	rootCmd.AddCommand(cmd)
}
