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

func response_headers_Render(m *Model) string {
	var content string

	keyTextStyle := lipgloss.NewStyle().Foreground(theme.Primary)

	keys := make([]string, 0, len(m.response.Headers))
	for key, value := range m.response.Headers {
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	for i, key := range keys {
		value := m.response.Headers[key]
		content += ansi.Wrap(keyTextStyle.Render(key)+": "+value[0], term_size.GetTerminalWidth(9999), "")
		if i < len(keys)-1 {
			content += "\n"
		}
	}
	return content
}

func response_headers_Render_Paged(m *Model) string {
	content := response_headers_Render(m)
	lines := strings.Split(content, "\n")

	m.respHeadersLinesAmount = len(lines)

	maxLines := term_size.GetTerminalHeight(9999) - 16
	start := min(m.respHeadersScrollOffset, len(lines))
	end := min(start+maxLines, len(lines))

	result := strings.Join(lines[start:end], "\n")

	if len(lines) > maxLines {
		keyTextStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
		result += keyTextStyle.Render(fmt.Sprintf("\n[%d-%d/%d lines] Use ↑/↓ or PgUp/PgDown to scroll.", start+1, end, len(lines)))
	}

	return result
}

func response_headers_ScrollUp(m *Model) {
	if m.respHeadersScrollOffset > 0 {
		m.respHeadersScrollOffset--
	}
}

func response_headers_ScrollDown(m *Model) {
	maxLines := term_size.GetTerminalHeight(9999) - 16

	if m.respHeadersLinesAmount == 0 {
		return
	}
	if m.respHeadersLinesAmount <= maxLines {
		return
	}

	if m.respHeadersScrollOffset+maxLines >= m.respHeadersLinesAmount {
		return
	} else {
		m.respHeadersScrollOffset++
	}
}

func response_headers_ScrollPgUp(m *Model) {
	m.respHeadersScrollOffset -= 5
	if m.respHeadersScrollOffset < 0 {
		m.respHeadersScrollOffset = 0
	}
}

func response_headers_ScrollPgDown(m *Model) {
	maxLines := term_size.GetTerminalHeight(9999) - 16
	if m.respHeadersLinesAmount == 0 || m.respHeadersLinesAmount <= maxLines {
		return
	}

	m.respHeadersScrollOffset += 5
}
