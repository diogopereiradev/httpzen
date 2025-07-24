package request_menu

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

func request_headers_Render(m *Model) string {
	var content string

	keyTextStyle := lipgloss.NewStyle().Foreground(theme.Primary)

	if len(m.response.Request.Headers) < 1 {
		return lipgloss.NewStyle().Foreground(theme.Warn).Render("No request headers found.")
	}

	keys := make([]string, 0, len(m.response.Request.Headers))
	for key, value := range m.response.Request.Headers {
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	for i, key := range keys {
		value := m.response.Request.Headers[key]
		content += ansi.Wrap(keyTextStyle.Render(key)+": "+value[0], term_size.GetTerminalWidth(9999), "")
		if i < len(keys)-1 {
			content += "\n"
		}
	}
	return content
}

func request_headers_Render_Paged(m *Model) string {
	content := request_headers_Render(m)
	lines := strings.Split(content, "\n")

	m.reqHeadersLinesAmount = len(lines)

	maxLines := term_size.GetTerminalHeight(9999) - 16
	start := min(m.reqHeadersScrollOffset, len(lines))
	end := min(start+maxLines, len(lines))

	result := strings.Join(lines[start:end], "\n")

	if len(lines) > maxLines {
		keyTextStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
		result += keyTextStyle.Render(fmt.Sprintf("\n[%d-%d/%d lines] Use ↑/↓ or PgUp/PgDown to scroll.", start+1, end, len(lines)))
	}

	return result
}

func request_headers_ScrollUp(m *Model) {
	if m.reqHeadersScrollOffset > 0 {
		m.reqHeadersScrollOffset--
	}
}

func request_headers_ScrollDown(m *Model) {
	maxLines := term_size.GetTerminalHeight(9999) - 16

	if m.reqHeadersLinesAmount == 0 {
		return
	}
	if m.reqHeadersLinesAmount <= maxLines {
		return
	}

	if m.reqHeadersScrollOffset+maxLines >= m.reqHeadersLinesAmount {
		return
	} else {
		m.reqHeadersScrollOffset++
	}
}

func request_headers_ScrollPgUp(m *Model) {
	m.reqHeadersScrollOffset -= 5
	if m.reqHeadersScrollOffset < 0 {
		m.reqHeadersScrollOffset = 0
	}
}

func request_headers_ScrollPgDown(m *Model) {
	maxLines := term_size.GetTerminalHeight(9999) - 16
	if m.reqHeadersLinesAmount == 0 || m.reqHeadersLinesAmount <= maxLines {
		return
	}

	m.reqHeadersScrollOffset += 5
}
