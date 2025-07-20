package logger_module

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Error(message string) {
	borderStyle := lipgloss.
		NewStyle().
		Width(80).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9")).
		Padding(1, 1)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

	var content string
	content += titleStyle.Render("An error has occurred on application execution") + "\n\n"
	content += messageStyle.Render(message)

	fmt.Println(borderStyle.Render(content))
}

func Warn(message string) {
	borderStyle := lipgloss.
		NewStyle().
		Width(50).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("3")).
		Padding(1, 1)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

	var content string
	content += titleStyle.Render("A warning has occurred") + "\n\n"
	content += messageStyle.Render(message)

	fmt.Println(borderStyle.Render(content))
}

func Info(message string) {
	borderStyle := lipgloss.
		NewStyle().
		Width(50).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		Padding(1, 1)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

	var content string
	content += titleStyle.Render("Information") + "\n\n"
	content += messageStyle.Render(message)

	fmt.Println(borderStyle.Render(content))
}
