package help_command

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		borderStyle := lipgloss.
			NewStyle().
			Width(85).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("6")).
			Padding(0, 3)

		logoStyle := lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("6")).
			Bold(true)

		titleStyle := lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("3")).
			Bold(true).
			Underline(true)

		fieldStyle := lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("#575757ff")).
			Bold(true)

		var content string
		content += logoStyle.Render(logoascii.GetLogo()) + "\n"
		content += titleStyle.Render("Httpzen CLI Tool for API Management and Development") + "\n"
		content += "\n"
		content += "Usage\n"
		content += "  $ httpzen " + fieldStyle.Render("[METHOD] ") + fieldStyle.Render("[URL] ") + fieldStyle.Render("[PARAMETERS ...]")

		if len(cmd.Commands()) > 0 {
			content += "\n\n"
			content += "Main Commands\n\n"
		}

		maxUseLen := 0
		for _, c := range cmd.Commands() {
			if c.Use == "completion" || strings.HasPrefix(c.Use, "help") {
				continue
			}
			if len(c.Use) > maxUseLen {
				maxUseLen = len(c.Use)
			}
		}

		for _, cmd := range cmd.Commands() {
			if cmd.Use == "completion" || strings.HasPrefix(cmd.Use, "help") {
				continue
			}
			padding := strings.Repeat(" ", maxUseLen-len(cmd.Use)+4)
			content += "  " + cmd.Use + padding + cmd.Short + "\n"
		}

		if cmd.Flags().HasFlags() {
			content += "\n"
			content += "Available parameters\n\n"
		}

		for _, flag := range cmd.Flags().FlagUsages() {
			flagStr := string(flag)
			content += flagStr
		}

		fmt.Println(borderStyle.Render(content))
	})
}
