package keyvalue_menu_component

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

var New = NewComponent

type KeyValue struct {
	Key   string
	Value string
}

type KeyValueMenuImpl struct {
	Title    string
	OnSubmit func([]KeyValue)
}

type model struct {
	title    string
	max      int
	pairs    []KeyValue
	adding   bool
	selected int
	keyInput string
	valInput string
	phase    string // "key", "value", "done"
	onSubmit func([]KeyValue)
}

func NewComponent(menu KeyValueMenuImpl) {
	m := &model{
		title:    menu.Title,
		onSubmit: menu.OnSubmit,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error on executing the menu:", err)
	}
}

func (m *model) Init() tea.Cmd {
	m.adding = true
	m.selected = 0
	m.phase = "key"
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.phase = "done"
			if m.onSubmit != nil {
				m.onSubmit(m.pairs)
				m.onSubmit = nil
			}
			return m, tea.Quit
		case tea.KeyEnter:
			switch m.phase {
			case "key":
				if m.keyInput == "" {
					m.phase = "done"
					if m.onSubmit != nil {
						m.onSubmit(m.pairs)
						m.onSubmit = nil
					}
					return m, tea.Quit
				}
				m.phase = "value"
				return m, nil
			case "value":
				m.pairs = append(m.pairs, KeyValue{Key: m.keyInput, Value: m.valInput})
				m.keyInput = ""
				m.valInput = ""
				if m.max > 0 && len(m.pairs) >= m.max {
					m.phase = "done"
					if m.onSubmit != nil {
						m.onSubmit(m.pairs)
						m.onSubmit = nil
					}
					return m, tea.Quit
				}
				m.phase = "key"
				return m, nil
			}
		case tea.KeyBackspace, tea.KeyCtrlH:
			if m.phase == "key" && len(m.keyInput) > 0 {
				m.keyInput = m.keyInput[:len(m.keyInput)-1]
			} else if m.phase == "value" && len(m.valInput) > 0 {
				m.valInput = m.valInput[:len(m.valInput)-1]
			}
		case tea.KeyUp:
			if len(m.pairs) > 0 && m.phase == "key" {
				if m.selected > 0 {
					m.selected--
				}
			}
		case tea.KeyDown:
			if len(m.pairs) > 0 && m.phase == "key" {
				if m.selected < len(m.pairs)-1 {
					m.selected++
				}
			}
		case tea.KeyDelete:
			if len(m.pairs) > 0 && m.phase == "key" && m.selected >= 0 && m.selected < len(m.pairs) {
				m.pairs = append(m.pairs[:m.selected], m.pairs[m.selected+1:]...)
				if m.selected > 0 && m.selected >= len(m.pairs) {
					m.selected--
				}
			}
		}
		if m.phase == "key" && msg.Type == tea.KeyRunes {
			m.keyInput += msg.String()
		} else if m.phase == "value" && msg.Type == tea.KeyRunes {
			m.valInput += msg.String()
		}
		return m, nil
	}
	return m, nil
}

func (m *model) View() string {
	config := config_module.GetConfig()

	var output string

	styleLabel := lipgloss.NewStyle().Foreground(theme.Secondary).Bold(true)
	styleInput := lipgloss.NewStyle().Foreground(theme.LightText).Background(theme.DarkenText).Padding(0, 1)
	styleSelected := lipgloss.NewStyle().Background(theme.Primary).Foreground(theme.LightText).Bold(true)
	greyText := lipgloss.NewStyle().Foreground(theme.DarkenText)

	questionIcon := lipgloss.NewStyle().Foreground(theme.Warn).Render("?")
	checkIcon := lipgloss.NewStyle().Foreground(theme.Success).Render("✓")

	if !config.HideLogomark {
		output += lipgloss.NewStyle().Foreground(theme.Primary).Render(logoascii.GetLogo(".body-builder")) + "\n"
	}

	if m.phase == "done" {
		output += checkIcon + " " + m.title + "\n"
		return output
	}

	output += questionIcon + " " + m.title + "\n"

	if len(m.pairs) > 0 {
		t := table.New()
		t.Border(lipgloss.RoundedBorder())
		t.BorderStyle(lipgloss.NewStyle().Foreground(theme.Primary))
		t.Width(terminal_utility.GetTerminalWidth(80))

		t.Headers(greyText.Render("Key"), greyText.Render("Value"))
		for i, kv := range m.pairs {
			keyCell := kv.Key
			valCell := kv.Value
			if m.phase == "key" && i == m.selected {
				keyCell = styleSelected.Render(keyCell)
				valCell = styleSelected.Render(valCell)
			}
			t.Row(keyCell, valCell)
		}
		output += t.Render() + "\n\n"
	}

	switch m.phase {
	case "key":
		output += styleLabel.Render("Key:") + " " + styleInput.Render(m.keyInput) + "\n"
		output += lipgloss.NewStyle().Faint(true).Render("(Enter to confirm, empty to finish)  ↑/↓ navigate, DEL to remove")
	case "value":
		output += styleLabel.Render("Value to '") + m.keyInput + "': " + styleInput.Render(m.valInput) + "\n"
		output += lipgloss.NewStyle().Faint(true).Render("(Enter to add, ESC to cancel)")
	}

	return output
}
