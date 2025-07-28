package number_prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

var New = newComponent

type NumberPromptImpl struct {
	Title                 string
	IncorrectTitleMessage string
	Events                NumberPromptEvents
	Default               int
	MaxLength             int
	input                 textinput.Model
	invalid               bool
	submited              bool
}

type NumberPromptEvents struct {
	OnSubmit      func(result string)
	OnSubmitError func()
}

func newComponent(options NumberPromptImpl) {
	input := textinput.New()
	input.CharLimit = options.MaxLength
	input.Focus()
	input.Prompt = ""

	impl := NumberPromptImpl{
		Title:                 options.Title,
		IncorrectTitleMessage: options.IncorrectTitleMessage,
		Events:                options.Events,
		MaxLength:             options.MaxLength,
		Default:               options.Default,
		input:                 input,
		invalid:               false,
		submited:              false,
	}

	terminal_utility.Clear()

	p := tea.NewProgram(impl)
	if model, err := p.Run(); err != nil {
		fmt.Println("Error on executing the program:", err)
		os.Exit(1)
	} else {
		p := model.(NumberPromptImpl)
		if p.submited {
			if p.Events.OnSubmit != nil {
				p.Events.OnSubmit(strings.ToLower(p.input.Value()))
			}
		} else {
			exitStyle := lipgloss.NewStyle().Foreground(theme.Error).Bold(true)
			fmt.Println(exitStyle.Render("No value provided, program exited."))
			os.Exit(1)
		}
	}
}

func (p NumberPromptImpl) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
	)
}

func (p NumberPromptImpl) View() string {
	config := config_module.GetConfig()

	var output string

	questionIcon := lipgloss.NewStyle().Foreground(theme.Warn).Render("?")
	checkIcon := lipgloss.NewStyle().Foreground(theme.Success).Render("✓")
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Secondary)

	if p.invalid {
		invalidStyle := lipgloss.NewStyle().Foreground(theme.Error).Bold(true)
		return questionIcon + " " + p.Title + ": " + p.input.View() + "\n" + invalidStyle.Render(p.IncorrectTitleMessage)
	}

	if !config.HideLogomark {
		output += lipgloss.NewStyle().Foreground(theme.Primary).Render(logoascii.GetLogo(".number-prompt")) + "\n"
	}

	var defaultValueStr string
	if p.Default != 0 {
		defaultValueStr = lipgloss.NewStyle().Foreground(theme.DarkenText).Render(" (default: " + fmt.Sprint(p.Default) + ")")
	}

	if p.submited {
		output += checkIcon + " " + p.Title + defaultValueStr + ": " + selectedStyle.Render(p.input.Value()) + "\n"
	}
	output += questionIcon + " " + p.Title + defaultValueStr + ": " + p.input.View() + "\n"
	return output
}

func (p NumberPromptImpl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if p.submited {
		return p, tea.Quit
	}

	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)

	filtered := ""
	for _, r := range p.input.Value() {
		if r >= '0' && r <= '9' {
			filtered += string(r)
		}
	}
	if filtered != p.input.Value() {
		p.input.SetValue(filtered)
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyRunes:
			p.invalid = false
		case tea.KeyEnter:
			if p.input.Value() == "" && p.Default == 0 {
				p.invalid = true
				return p, nil
			}
			p.submited = true
			return p, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return p, tea.Quit
		}
	}
	return p, cmd
}
