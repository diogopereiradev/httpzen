package logger_module

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

func getLogWidth(maxWidth int) int {
	var width int

	if maxWidth > 0 {
		width = maxWidth
	} else {
		width = terminal_utility.GetTerminalWidth(90)
	}
	return width
}

func Error(message string, maxWidth int) {
	config := config_module.GetConfig()
	width := getLogWidth(maxWidth)

	borderStyle := lipgloss.
		NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9")).
		Padding(0, 2, 1, 2)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

	var content string
	if !config.HideLogomark {
		content += titleStyle.Render(logoascii.GetLogo(".logger")) + "\n"
	}
	content += titleStyle.Render("An error has occurred on application execution") + "\n\n"
	content += messageStyle.Render(message)

	fmt.Println(borderStyle.Render(content))
}

func Warn(message string, maxWidth int) {
	config := config_module.GetConfig()
	width := getLogWidth(maxWidth)

	borderStyle := lipgloss.
		NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("3")).
		Padding(0, 2, 1, 2)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

	var content string
	if !config.HideLogomark {
		content += titleStyle.Render(logoascii.GetLogo(".logger")) + "\n"
	}
	content += titleStyle.Render("A warning has occurred") + "\n\n"
	content += messageStyle.Render(message)

	fmt.Println(borderStyle.Render(content))
}

func Info(message string, maxWidth int) {
	config := config_module.GetConfig()
	width := getLogWidth(maxWidth)

	borderStyle := lipgloss.
		NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		Padding(0, 2, 1, 2)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

	var content string
	if !config.HideLogomark {
		content += titleStyle.Render(logoascii.GetLogo(".logger")) + "\n"
	}
	content += titleStyle.Render("Information") + "\n\n"
	content += messageStyle.Render(message)

	fmt.Println(borderStyle.Render(content))
}

func Success(message string, maxWidth int) {
	config := config_module.GetConfig()
	width := getLogWidth(maxWidth)

	var content string

	borderStyle := lipgloss.
		NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Success).
		Padding(0, 2, 1, 2)

	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.Success)

	if !config.HideLogomark {
		content += titleStyle.Render(logoascii.GetLogo(".logger")) + "\n"
	}

	content += titleStyle.Render("Action Succeeded") + "\n\n"
	content += messageStyle.Render(message)

	fmt.Println(borderStyle.Render(content))
}
