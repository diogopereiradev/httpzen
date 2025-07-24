package request_menu

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

type tab int

const (
	tab_Result tab = iota
	tab_RequestInfos
	tab_NetworkInfos
	tab_RequestHeaders
	tab_ResponseHeaders
)

var tabNames = []string{
	"Response",
	"Request Infos",
	"Network Infos",
	"Request Headers",
	"Response Headers",
}

var activeTabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      " ",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┘",
	BottomRight: "└",
}

var tabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┴",
	BottomRight: "┴",
}

func tab_Render(m *Model) string {
	tabStyle := lipgloss.NewStyle().
		BorderForeground(theme.Primary).
		Padding(0, 1)

	var tabLabels []string
	for i, name := range tabNames {
		if tab(i) == m.activeTab {
			tabLabels = append(tabLabels, tabStyle.Border(activeTabBorder).Foreground(theme.Primary).Render(name))
		} else {
			tabLabels = append(tabLabels, tabStyle.Border(tabBorder).Render(name))
		}
	}

	tabGap := tabStyle.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabLabels...,
	)

	gap := tabGap.Render(strings.Repeat(" ", max(0, term_size.GetTerminalWidth(9999)-lipgloss.Width(row))))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	return row
}

func tab_MoveLeft(m *Model) Model {
	m.activeTab = tab((int(m.activeTab) - 1 + len(tabNames)) % len(tabNames))
	m.resultScrollOffset = 0
	return *m
}

func tab_MoveRight(m *Model) Model {
	m.activeTab = tab((int(m.activeTab) + 1) % len(tabNames))
	m.resultScrollOffset = 0
	return *m
}
