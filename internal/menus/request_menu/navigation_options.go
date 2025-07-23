package request_menu

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

func navigation_options_Render() string {
	var content string
	greyTextStyle := lipgloss.NewStyle().Foreground(theme.DarkenText)

	content += greyTextStyle.Render("\n\nUse left/right arrows to navigate between tabs, 'q' to quit.")
	content += greyTextStyle.Render("\n'c' to copy response, 'b' to benchmark, and 'r' to resend request.\n")

	return content
}