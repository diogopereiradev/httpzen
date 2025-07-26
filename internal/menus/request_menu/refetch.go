package request_menu

import (
	"github.com/charmbracelet/lipgloss"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

type RefetchEvent struct {
	Response request_module.RequestResponse
}

func refetch_Render(m *Model) string {
	content := lipgloss.
		NewStyle().
		Width(terminal_utility.GetTerminalWidth(9999)).
		Foreground(theme.Primary).
		Render(logoascii.GetLogo(".request"))

	content += "\n\n"
	content += lipgloss.
		NewStyle().
		Width(terminal_utility.GetTerminalWidth(9999)).
		Foreground(theme.LightText).
		Background(theme.Primary).
		Align(lipgloss.Center).
		Render("Refetching request...")

	return content
}
