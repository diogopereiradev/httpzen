package request_menu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
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
	content += fieldTextStyle.Render("Response Time: ") + executionTime + "\n"
	content += fieldTextStyle.Render("Response Size: ") + fmt.Sprintf("%d bytes", len(m.response.Result))

	if len(m.response.Body) > 0 {
		content += "\n"
		content += fieldTextStyle.Render("\nRequest Body:\n" + m.response.Body[0].Value + "\n")
	}
	return content
}

func basic_infos_Render_Paged(m *Model) string {
	content := basic_infos_Render(m)
	lines := strings.Split(content, "\n")

	m.infosLinesAmount = len(lines)

	maxLines := term_size.GetTerminalHeight(9999) - 16
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
	maxLines := term_size.GetTerminalHeight(9999) - 16

	if m.infosLinesAmount == 0 {
		return
	}
	if m.infosLinesAmount <= maxLines {
		return
	}

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
	maxLines := term_size.GetTerminalHeight(9999) - 16
	if m.infosLinesAmount == 0 || m.infosLinesAmount <= maxLines {
		return
	}

	m.infosScrollOffset += 5
}
