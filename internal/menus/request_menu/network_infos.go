package request_menu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

func network_infos_Render(m *Model) string {
	var content string

	fieldTextStyle := lipgloss.NewStyle().Foreground(theme.Primary)

	var blocks []string
	for _, info := range m.response.IpInfos {
		var block string
		block += fieldTextStyle.Render("Protocol: ") + info.Type + "\n"
		block += fieldTextStyle.Render("IP Address: ") + info.Ip + "\n"

		if info.Country != "" {
			block += fieldTextStyle.Render("Country: ") + info.Country + "\n"
		}
		if info.City != "" {
			block += fieldTextStyle.Render("City: ") + info.City + "\n"
		}
		if info.Decimal != "" {
			block += fieldTextStyle.Render("Decimal: ") + info.Decimal + "\n"
		}
		if info.Hostname != "" {
			block += fieldTextStyle.Render("Hostname: ") + info.Hostname + "\n"
		}
		if info.State != "" {
			block += fieldTextStyle.Render("Region/State: ") + info.State + "\n"
		}
		if info.ASN != "" {
			block += fieldTextStyle.Render("ASN: ") + info.ASN + "\n"
		}
		if info.ISP != "" {
			block += fieldTextStyle.Render("ISP: ") + info.ISP + "\n"
		}

		if info.Latitude != 0 && info.Longitude != 0 {
			block += fieldTextStyle.Render("Coordinates: ") + fmt.Sprintf("%.2f, %.2f", info.Latitude, info.Longitude)
		}
		blocks = append(blocks, block)
	}

	content += strings.TrimRight(strings.Join(blocks, "\n\n"), "\n ")
	return content
}

func network_infos_Render_Paged(m *Model) string {
	content := network_infos_Render(m)
	lines := strings.Split(content, "\n")

	m.networkLinesAmount = len(lines)

	maxLines := term_size.GetTerminalHeight(9999) - 16

	start := min(m.networkScrollOffset, len(lines))
	end := min(start+maxLines, len(lines))

	result := strings.Join(lines[start:end], "\n")

	if len(lines) > maxLines {
		fieldTextStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
		result += fieldTextStyle.Render(fmt.Sprintf("\n[%d-%d/%d lines] Use ↑/↓ or PgUp/PgDown to scroll.", start+1, end, len(lines)))
	}

	return result
}

func network_infos_ScrollUp(m *Model) {
	if m.networkScrollOffset > 0 {
		m.networkScrollOffset--
	}
}

func network_infos_ScrollDown(m *Model) {
	maxLines := term_size.GetTerminalHeight(9999) - 16

	if m.networkLinesAmount == 0 {
		return
	}
	if m.networkLinesAmount <= maxLines {
		return
	}

	if m.networkScrollOffset+maxLines >= m.networkLinesAmount {
		return
	} else {
		m.networkScrollOffset++
	}
}

func network_infos_ScrollPgUp(m *Model) {
	m.networkScrollOffset -= 5
	if m.networkScrollOffset < 0 {
		m.networkScrollOffset = 0
	}
}

func network_infos_ScrollPgDown(m *Model) {
	maxLines := term_size.GetTerminalHeight(9999) - 16
	if m.networkLinesAmount == 0 || m.networkLinesAmount <= maxLines {
		return
	}

	m.networkScrollOffset += 5
}
