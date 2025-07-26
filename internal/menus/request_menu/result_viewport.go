package request_menu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/diogopereiradev/httpzen/internal/utils/html_formatter"
	"github.com/diogopereiradev/httpzen/internal/utils/http_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/json_formatter"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

func result_viewport_Render(m *Model) string {
	var result string

	fieldTextStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
	maxLines := terminal_utility.GetTerminalHeight(9999) - 16

	if m.response.Result != "" {
		var formatted string
		contentType := http_utility.DetectContentType(m.response.Result)

		switch contentType {
		case "json":
			formatted = ansi.Wrap(json_formatter.FormatJSON(m.response.Result), terminal_utility.GetTerminalWidth(9999), "")
		case "html":
			formatted = ansi.Wrap(html_formatter.FormatHTML(m.response.Result), terminal_utility.GetTerminalWidth(9999), "")
		case "xml":
			formatted = ansi.Wrap(html_formatter.FormatHTML(m.response.Result), terminal_utility.GetTerminalWidth(9999), "")
		default:
			formatted = ansi.Wrap(m.response.Result, terminal_utility.GetTerminalWidth(9999), "")
		}

		lines := strings.Split(formatted, "\n")
		total := len(lines)

		if m.resultScrollOffset > total-maxLines {
			m.resultScrollOffset = max(0, total-maxLines)
		}

		m.resultLinesAmount = total

		end := min(m.resultScrollOffset+maxLines, total)
		visible := lines[m.resultScrollOffset:end]
		result += strings.Join(visible, "\n")

		if total > maxLines {
			result += fieldTextStyle.Render(fmt.Sprintf("\n[%d-%d/%d lines] Use ↑/↓ or PgUp/PgDown to scroll.", m.resultScrollOffset+1, end, total))
		}
	} else {
		result += fieldTextStyle.Render("No response available.")
	}
	return result
}

func result_viewport_ScrollUp(m *Model) {
	if m.resultScrollOffset > 0 {
		m.resultScrollOffset--
	}
}

func result_viewport_ScrollDown(m *Model) {
	if m.resultLinesAmount == 0 {
		return
	}
	if m.resultScrollOffset >= m.resultLinesAmount {
		m.resultScrollOffset = m.resultLinesAmount - 1
	} else {
		m.resultScrollOffset++
	}
}

func result_viewport_ScrollPgUp(m *Model) {
	m.resultScrollOffset -= 5
	if m.resultScrollOffset < 0 {
		m.resultScrollOffset = 0
	}
}

func result_viewport_ScrollPgDown(m *Model) {
	if m.resultLinesAmount == 0 || m.resultLinesAmount <= terminal_utility.GetTerminalHeight(9999)-16 {
		return
	}
	m.resultScrollOffset += 5
}
