package textarea_component

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

var New = NewComponent

type TextareaImpl struct {
	Title     string
	Events    TextareaEvents
	MaxLength int
	input     textarea.Model
	submited  bool
}

type TextareaEvents struct {
	OnSubmit func(result string)
}

func NewComponent(options TextareaImpl) {
	input := textarea.New()
	input.CharLimit = options.MaxLength
	input.FocusedStyle.LineNumber = lipgloss.NewStyle().Foreground(theme.DarkenText)
	input.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(theme.DarkenText)
	input.Focus()
	input.Prompt = ""
	input.ShowLineNumbers = true
	input.SetWidth(60)

	impl := TextareaImpl{
		Title:     options.Title,
		Events:    options.Events,
		MaxLength: options.MaxLength,
		input:     input,
		submited:  false,
	}

	p := tea.NewProgram(impl)
	if model, err := p.Run(); err != nil {
		fmt.Println("Error on executing the program:", err)
		os.Exit(1)
	} else {
		p := model.(TextareaImpl)
		if p.submited {
			if p.Events.OnSubmit != nil {
				p.Events.OnSubmit(p.input.Value())
			}
		} else {
			exitStyle := lipgloss.NewStyle().Foreground(theme.Error).Bold(true)
			fmt.Println(exitStyle.Render("No value provided, program exited."))
			os.Exit(1)
		}
	}
}

func (p TextareaImpl) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
	)
}

func (p TextareaImpl) View() string {
	config := config_module.GetConfig()

	var output string

	questionIcon := lipgloss.NewStyle().Foreground(theme.Warn).Render("?")
	checkIcon := lipgloss.NewStyle().Foreground(theme.Success).Render("✓")
	optionalStyle := lipgloss.NewStyle().Foreground(theme.DarkenText)
	borderStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(80).
		UnsetBackground()

	vlen := fmt.Sprintf("%d", len(p.input.Value()))
	if len(p.input.Value()) < 10 {
		vlen = "0" + vlen
	}

	maxlen := fmt.Sprintf("%d", p.MaxLength)
	if p.MaxLength < 10 {
		maxlen = "0" + maxlen
	}

	if !config.HideLogomark {
		output += lipgloss.NewStyle().Foreground(theme.Primary).Render(logoascii.GetLogo(".body-builder")) + "\n"
	}

	if p.submited {
		submitedStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
		output = checkIcon + " " + p.Title + submitedStyle.Render(" ["+vlen+"/"+maxlen+"]") + "\n"
		return output
	}
	output += questionIcon + " " + p.Title + " [" + vlen + "/" + maxlen + "]" + optionalStyle.Render("(optional)") + ": " + "\n\n"
	output += optionalStyle.Render("(You can use \"CTRL + S\" to submit, \"Enter\" to add a new line)") + "\n"
	output += borderStyle.Render(p.input.View()) + "\n"

	return output
}

func (p TextareaImpl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if p.submited {
		return p, tea.Quit
	}

	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyCtrlS:
			p.submited = true
			return p, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return p, tea.Quit
		}
	}
	return p, cmd
}