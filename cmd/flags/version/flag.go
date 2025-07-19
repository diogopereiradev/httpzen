package version_flag

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var Version string = "unknown"
var BuildDate string = "unknown"
var Website string = "unknown"
var Repository string = "unknown"
var License string = "unknown"

var Exit = os.Exit

func AddFlag(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show the version of the application")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("227"))

		borderStyle := lipgloss.
			NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Padding(1, 2)

		fieldKeyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))

		var versionInfo string
		versionInfo += titleStyle.Render("HTTPZen - Version information") + "\n\n"
		versionInfo += fieldKeyStyle.Render("Version: ") + Version + "\n"
		versionInfo += fieldKeyStyle.Render("Build Date: ") + BuildDate + "\n"
		versionInfo += fieldKeyStyle.Render("Website: ") + Website + "\n"
		versionInfo += fieldKeyStyle.Render("Repository: ") + Repository + "\n"
		versionInfo += fieldKeyStyle.Render("License: ") + License

		fmt.Println(borderStyle.Render(versionInfo))

		Exit(0)
	}
}
