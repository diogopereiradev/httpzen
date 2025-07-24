package version_command

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
	"github.com/spf13/cobra"
)

var Version string = "unknown"
var BuildDate string = "unknown"
var Website string = "unknown"
var Repository string = "unknown"
var License string = "unknown"

var Exit = os.Exit

func Init(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version of the application",
		Run: func(cmd *cobra.Command, args []string) {
			titleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.Primary)

			borderStyle := lipgloss.
				NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Primary).
				Padding(1, 2)

			fieldKeyStyle := lipgloss.NewStyle().Foreground(theme.Secondary)

			var versionInfo string
			versionInfo += titleStyle.Render("HTTPZen - Version information") + "\n\n"
			versionInfo += fieldKeyStyle.Render("Version: ") + Version + "\n"
			versionInfo += fieldKeyStyle.Render("Build Date: ") + BuildDate + "\n"
			versionInfo += fieldKeyStyle.Render("Website: ") + Website + "\n"
			versionInfo += fieldKeyStyle.Render("Repository: ") + Repository + "\n"
			versionInfo += fieldKeyStyle.Render("License: ") + License

			fmt.Println(borderStyle.Render(versionInfo))

			Exit(0)
		},
	}
	rootCmd.AddCommand(cmd)
}
