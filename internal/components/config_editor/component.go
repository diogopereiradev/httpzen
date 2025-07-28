package config_editor

import (
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	blink_cursor_component "github.com/diogopereiradev/httpzen/internal/components/blink_cursor"
	timed_message_component "github.com/diogopereiradev/httpzen/internal/components/timed_message"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

type configType int

const (
	optionTypeNumber configType = iota
	optionTypeBool
)

type configOption struct {
	Type      configType
	ConfigKey string
	Label     string
	Value     any
}

type model struct {
	config  config_module.Config
	options []configOption

	choice int

	editing      bool
	editingValue string

	savedTimedMessage *timed_message_component.TimedMessage
	blinkCursor       *blink_cursor_component.BlinkCursor
}

func getConfigOptions(config *config_module.Config) []configOption {
	return []configOption{
		{
			Type:      optionTypeNumber,
			ConfigKey: "SlowResponseThreshold",
			Label:     "Slow response threshold(ms)",
			Value:     config.SlowResponseThreshold,
		},
		{
			Type:      optionTypeBool,
			ConfigKey: "HideLogomark",
			Label:     "Hide logomark",
			Value:     config.HideLogomark,
		},
	}
}

func initialModel() model {
	config := config_module.GetConfig()
	return model{
		config:            config,
		options:           getConfigOptions(&config),
		savedTimedMessage: timed_message_component.New(),
		blinkCursor:       blink_cursor_component.New(),
		choice:            0,
		editing:           false,
	}
}

func RunConfigEditor() error {
	p := tea.NewProgram(initialModel())
	terminal_utility.Clear()

	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	var view string

	titleStyle := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true)
	cursorStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
	fieldStyle := lipgloss.NewStyle().Foreground(theme.Success)
	greyStyle := lipgloss.NewStyle().Foreground(theme.DarkenText)
	selectedShortcutStyle := lipgloss.NewStyle().Foreground(theme.Secondary).Bold(true)
	valueStyle := lipgloss.NewStyle().Background(theme.CodeBlock).Foreground(theme.LightText).Padding(0, 1)

	borderStyle := lipgloss.
		NewStyle().
		Width(terminal_utility.GetTerminalWidth(90)-2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		Padding(0, 3, 1, 3)

	if !m.config.HideLogomark {
		view += titleStyle.Render(logoascii.GetLogo(".config")) + "\n"
	} else {
		borderStyle = borderStyle.Padding(1, 3, 1, 3)
	}

	cursor := func(i int) string {
		if m.choice == i {
			return cursorStyle.Render("> ")
		}
		return " "
	}

	view += titleStyle.Render("Configuration:") + "\n\n"

	for i, opt := range m.options {
		view += fmt.Sprintf("%s "+fieldStyle.Render(opt.Label+":")+" %v\n", cursor(i), opt.Value)
	}

	blinkCursor := m.blinkCursor.Render()

	if m.editing {
		choice := m.options[m.choice]
		switch choice.Type {
		case optionTypeBool:
			view += "\n" + selectedShortcutStyle.Render("Press 'enter' to toggle('n' to cancel)")
		case optionTypeNumber:
			view += "\n" + selectedShortcutStyle.Render("Enter new value and press 'enter' to save('n' to cancel): ") + valueStyle.Render(m.editingValue+string(blinkCursor))
		}
	}

	view += greyStyle.Render("\nUse ↑/↓ to navigate, Enter to edit, q to quit.")
	view += "\n\n" + selectedShortcutStyle.Bold(false).Render("Config located at: ") + greyStyle.Render(config_module.GetConfigFilePath())
	view += "\n" + selectedShortcutStyle.Bold(false).Render("Documentation: ") + greyStyle.Render("https://httpzen.diogopereira.site/docs/configuration")

	if m.savedTimedMessage != nil && m.savedTimedMessage.Visible {
		message := m.savedTimedMessage.Render()
		if message != "" {
			view += "\n\n" + message
		}
	}
	return borderStyle.Render(view)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Messages
	switch msg.(type) {
	case blink_cursor_component.BlinkCursorMsg:
		if !m.editing {
			m.blinkCursor.Hide()
			return m, nil
		}
		return m, m.blinkCursor.Show()
	}

	// Shortcuts
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m.quit()
		case tea.KeyRunes:
			return m.keyRunes(&msg)
		case tea.KeyBackspace:
			return m.backspace()
		case tea.KeyUp:
			return m.navigateUp()
		case tea.KeyDown:
			return m.navigateDown()
		case tea.KeyEnter:
			if !m.editing {
				m.editing = true
				m.editingValue = fmt.Sprintf("%v", m.options[m.choice].Value)
				return m, m.blinkCursor.Show()
			} else {
				return m.saveChanges()
			}
		}
	}
	return m, nil
}

func (m *model) keyRunes(msg *tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m.quit()
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		choice := m.options[m.choice]
		if m.editing && choice.Type == optionTypeNumber {
			m.editingValue += msg.String()
		}
	case "y", "Y":
		return m.saveChanges()
	case "n", "N":
		m.editing = false
		m.editingValue = ""
		return m, nil
	}
	return m, nil
}

func (m *model) quit() (tea.Model, tea.Cmd) {
	terminal_utility.Clear()
	return m, tea.Quit
}

func (m *model) navigateUp() (tea.Model, tea.Cmd) {
	if m.choice > 0 {
		m.choice--
	}
	return m, nil
}

func (m *model) navigateDown() (tea.Model, tea.Cmd) {
	if m.choice < len(m.options)-1 {
		m.choice++
	}
	return m, nil
}

func (m *model) backspace() (tea.Model, tea.Cmd) {
	choice := m.options[m.choice]

	if !m.editing || choice.Type != optionTypeNumber {
		return m, nil
	}
	if m.editing && len(m.editingValue) > 0 {
		m.editingValue = m.editingValue[:len(m.editingValue)-1]
	}
	return m, nil
}

func (m *model) saveChanges() (tea.Model, tea.Cmd) {
	if !m.editing {
		return m, nil
	}

	choice := m.options[m.choice]
	newConfig := m.config

	var setters = map[string]func(*config_module.Config){
		"HideLogomark": func(cfg *config_module.Config) { cfg.HideLogomark = !cfg.HideLogomark },
		"SlowResponseThreshold": func(cfg *config_module.Config) {
			if m.editingValue == "" {
				cfg.SlowResponseThreshold = 0
				return
			}
			if newValue, err := strconv.Atoi(m.editingValue); err == nil {
				cfg.SlowResponseThreshold = newValue
			}
		},
	}

	if setter, ok := setters[choice.ConfigKey]; ok {
		setter(&newConfig)
	}

	m.config = newConfig
	m.options = getConfigOptions(&newConfig)
	config_module.UpdateConfig(newConfig)

	m.editing = false
	m.editingValue = ""

	return m, m.savedTimedMessage.Show("Saved", 1*time.Second)
}
