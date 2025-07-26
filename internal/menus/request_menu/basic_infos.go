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

func basic_infos_Render(m *Model) string {
	var content string

	greyTextStyle := lipgloss.NewStyle().Foreground(theme.DarkenText)
	fieldTextStyle := lipgloss.NewStyle().Foreground(theme.Secondary)

	greenFieldTextStyle := lipgloss.NewStyle().Foreground(theme.Success)
	yellowFieldTextStyle := lipgloss.NewStyle().Foreground(theme.Warn)
	redFieldTextStyle := lipgloss.NewStyle().Foreground(theme.Error)

	var executionTime string
	if m.response.ExecutionTime > float64(m.config.SlowResponseThreshold) {
		executionTime = redFieldTextStyle.Render(fmt.Sprintf("%.2f", m.response.ExecutionTime) + "ms (slow)")
	} else if m.response.ExecutionTime > (float64(m.config.SlowResponseThreshold) * 0.7) {
		executionTime = yellowFieldTextStyle.Render(fmt.Sprintf("%.2f", m.response.ExecutionTime) + "ms (slow)")
	} else {
		executionTime = greenFieldTextStyle.Render(fmt.Sprintf("%.2f", m.response.ExecutionTime) + "ms (fast)")
	}

	content += greyTextStyle.Render(fmt.Sprint(m.response.HttpVersion)+" "+m.response.Method+" "+m.response.StatusMessage) + "\n\n"
	content += fieldTextStyle.Render("URL: ") + m.response.Request.Url + "\n"
	content += fieldTextStyle.Render("Response Time: ") + executionTime + "\n"
	content += fieldTextStyle.Render("Response Size: ") + fmt.Sprintf("%d bytes", len(m.response.Result))
	content += basic_infos_body_Render(m)

	return content
}

func basic_infos_Render_Paged(m *Model) string {
	content := basic_infos_Render(m)
	lines := strings.Split(content, "\n")

	m.infosLinesAmount = len(lines)

	maxLines := terminal_utility.GetTerminalHeight(9999) - 16
	start := min(m.infosScrollOffset, len(lines))
	end := min(start+maxLines, len(lines))

	result := strings.Join(lines[start:end], "\n")

	if len(lines) > maxLines {
		fieldTextStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
		result += fieldTextStyle.Render(fmt.Sprintf("\n[%d-%d/%d lines] Use ↑/↓ or PgUp/PgDown to scroll.", start+1, end, len(lines)))
	}

	return result
}

func basic_infos_ScrollUp(m *Model) {
	if m.infosScrollOffset > 0 {
		m.infosScrollOffset--
	}
}

func basic_infos_ScrollDown(m *Model) {
	maxLines := terminal_utility.GetTerminalHeight(9999) - 16

	if m.infosLinesAmount == 0 { return }
	if m.infosLinesAmount <= maxLines { return }

	if m.infosScrollOffset+maxLines >= m.infosLinesAmount {
		return
	} else {
		m.infosScrollOffset++
	}
}

func basic_infos_ScrollPgUp(m *Model) {
	m.infosScrollOffset -= 5
	if m.infosScrollOffset < 0 {
		m.infosScrollOffset = 0
	}
}

func basic_infos_ScrollPgDown(m *Model) {
	maxLines := terminal_utility.GetTerminalHeight(9999) - 16
	if m.infosLinesAmount == 0 || m.infosLinesAmount <= maxLines {
		return
	}

	m.infosScrollOffset += 5
}

func basic_infos_body_Render(m *Model) string {
	var content string

	greenFieldTextStyle := lipgloss.NewStyle().Foreground(theme.Success)
	yellowFieldTextStyle := lipgloss.NewStyle().Foreground(theme.Warn)
	redFieldTextStyle := lipgloss.NewStyle().Foreground(theme.Error)

	fieldTextStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
	specialFieldTextStyle := lipgloss.NewStyle().Foreground(theme.Primary)
	requestBodyStyle := lipgloss.NewStyle().Width(terminal_utility.GetTerminalWidth(9999)).Background(theme.CodeBlock).Padding(1, 1)

	content += "\n"
	content += fieldTextStyle.Render("Request Body:") + "\n\n"

	if len(m.response.Body) == 0 {
		content += yellowFieldTextStyle.Render("No request body available.") + "\n"
		return content;
	}

	for i, body := range m.response.Body {
		if i > 0 && i < len(m.response.Body) {
			content += "\n\n"
		}
		content += specialFieldTextStyle.Render("Content type: ") + body.ContentType + "\n"
		content += specialFieldTextStyle.Render("Content length: ") + fmt.Sprintf("%d bytes", len(body.Value)) + "\n"

		if fileInfo, err := http_utility.GetFileByPath(body.Value); err == nil {
			content += specialFieldTextStyle.Render("File status: ") + greenFieldTextStyle.Render("Found and accessible") + "\n"
			content += specialFieldTextStyle.Render("File name: ") + fileInfo.Name + "\n"
		} else if fileInfo.PathIsValid {
			content += specialFieldTextStyle.Render("File status: ") + redFieldTextStyle.Render("Not found or inaccessible") + "\n"
		}

		content += specialFieldTextStyle.Render("Content: ") + "\n\n"
		
		if body.Key != "" {
			content += ansi.Wrap(greenFieldTextStyle.Render("Key: ") + body.Key, terminal_utility.GetTerminalWidth(9999), "") + "\n"
			content += ansi.Wrap(greenFieldTextStyle.Render("Value: ") + body.Value, terminal_utility.GetTerminalWidth(9999), "")
		} else {
			switch body.ContentType {
		  case "application/json":
			  content += ansi.Wrap(requestBodyStyle.Render(json_formatter.FormatJSON(body.Value)), terminal_utility.GetTerminalWidth(9999), "") + "\n"
			case "text/plain":
				content += ansi.Wrap(requestBodyStyle.Render(body.Value), terminal_utility.GetTerminalWidth(9999), "") + "\n"
			case "text/html":
				content += ansi.Wrap(html_formatter.FormatHTML(requestBodyStyle.Render(body.Value)), terminal_utility.GetTerminalWidth(9999), "") + "\n"
			default:
				content += ansi.Wrap(requestBodyStyle.Render(body.Value), terminal_utility.GetTerminalWidth(9999), "") + "\n"
			}
		}

		if i >= len(m.response.Body) - 1 {
  		content += "\n"
 		}
	}
	return content
}
