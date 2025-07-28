package prompt

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

var New = NewComponent

type PromptImpl struct {
	Title                 string
	IncorrectTitleMessage string
	Events                PromptEvents
	Boolean               bool
	BooleanDefault        bool
	MaxLength             int
	input                 textinput.Model
	invalid               bool
	submited              bool
}

type PromptEvents struct {
	OnSubmit      func(result string)
	OnSubmitError func()
}

func NewComponent(options PromptImpl) {
	input := textinput.New()
	input.CharLimit = options.MaxLength
	input.Focus()
	input.Prompt = ""

	impl := PromptImpl{
		Title:                 options.Title,
		IncorrectTitleMessage: options.IncorrectTitleMessage,
		Events:                options.Events,
		Boolean:               options.Boolean,
		BooleanDefault:        options.BooleanDefault,
		MaxLength:             options.MaxLength,
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
		p := model.(PromptImpl)
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

func (p PromptImpl) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
	)
}

func (p PromptImpl) View() string {
	config := config_module.GetConfig()

	var output string

	questionIcon := lipgloss.NewStyle().Foreground(theme.Warn).Render("?")
	checkIcon := lipgloss.NewStyle().Foreground(theme.Success).Render("✓")
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Secondary)

	if p.Boolean {
		optionStyle := lipgloss.NewStyle().Foreground(theme.DarkenText)
		var option string
		if p.BooleanDefault {
			option = optionStyle.Render(" (Y/n)")
		} else {
			option = optionStyle.Render(" (y/N)")
		}

		if p.invalid {
			invalidStyle := lipgloss.NewStyle().Foreground(theme.Error).Bold(true)
			return questionIcon + " " + p.Title + option + ": " + p.input.View() + "\n" + invalidStyle.Render("Invalid input, please enter 'y' or 'n'.")
		}

		if p.submited {
			return checkIcon + " " + p.Title + option + ": " + selectedStyle.Render(p.input.Value()) + "\n"
		}
		return questionIcon + " " + p.Title + option + ": " + p.input.View() + "\n"
	}

	vlen := fmt.Sprintf("%d", len(p.input.Value()))
	if len(p.input.Value()) < 10 {
		vlen = "0" + vlen
	}

	maxlen := fmt.Sprintf("%d", p.MaxLength)
	if p.MaxLength < 10 {
		maxlen = "0" + maxlen
	}

	if p.invalid {
		invalidStyle := lipgloss.NewStyle().Foreground(theme.Error).Bold(true)
		return questionIcon + " " + p.Title + " [" + vlen + "/" + maxlen + "]" + ": " + p.input.View() + "\n" + invalidStyle.Render(p.IncorrectTitleMessage)
	}

	if !config.HideLogomark {
		output += lipgloss.NewStyle().Foreground(theme.Primary).Render(logoascii.GetLogo(".body-builder")) + "\n"
	}

	if p.submited {
		output += checkIcon + " " + p.Title + " [" + vlen + "/" + maxlen + "]" + ": " + selectedStyle.Render(p.input.Value()) + "\n"
	}
	output += questionIcon + " " + p.Title + " [" + vlen + "/" + maxlen + "]" + ": " + p.input.View() + "\n"
	return output
}

func (p PromptImpl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if p.submited {
		return p, tea.Quit
	}

	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyRunes:
			p.invalid = false
		case tea.KeyEnter:
			if p.Boolean {
				if p.input.Value() == "" {
					if p.BooleanDefault {
						p.input.SetValue("y")
					} else {
						p.input.SetValue("n")
					}
					p.submited = true
					return p, tea.Quit
				}
				if p.input.Value() == "y" || p.input.Value() == "Y" {
					p.submited = true
					return p, tea.Quit
				} else {
					if p.input.Value() == "n" || p.input.Value() == "N" {
						p.submited = true
						return p, tea.Quit
					}
					p.invalid = true
					return p, nil
				}
			}
			if p.input.Value() == "" {
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
