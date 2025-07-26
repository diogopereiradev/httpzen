package select_menu_component

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

type MenuImpl struct {
	Choices    []string
	Messages   MenuMessages
	PerPage    int
	Border     bool
	Events     MenuEvents
	cursor     int
	selected   int
	searchTerm string
}

type MenuMessages struct {
	Title        string
	EmptyOptions string
}

type MenuEvents struct {
	OnSelect func(choice int)
}

func New(options MenuImpl) {
	impl := MenuImpl{
		Choices:    options.Choices,
		Events:     options.Events,
		PerPage:    options.PerPage,
		Border:     options.Border,
		Messages:   options.Messages,
		cursor:     0,
		selected:   -1,
		searchTerm: "",
	}

	if options.PerPage <= 0 {
		options.PerPage = 10
	}

	view := impl.View()
	result, err := impl.Update(view)
	if err != nil {
		fmt.Println("Error on updating the view:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(result)
	if model, err := p.Run(); err != nil {
		fmt.Println("Error on executing the program:", err)
		os.Exit(1)
	} else {
		m := model.(MenuImpl)
		if m.selected != -1 {
			if m.Events.OnSelect != nil {
				m.Events.OnSelect(m.selected)
			}
		} else {
			exitStyle := lipgloss.NewStyle().Foreground(theme.Error).Bold(true)
			fmt.Println(exitStyle.Render("No option selected, program exited."))
			os.Exit(1)
		}
	}
}

func (m MenuImpl) FilteredChoices() []string {
	if m.searchTerm != "" {
		var filtered []string
		for _, choice := range m.Choices {
			if strings.Contains(strings.ToLower(choice), strings.ToLower(m.searchTerm)) {
				filtered = append(filtered, choice)
			}
		}
		return filtered
	}
	return m.Choices
}

func (m MenuImpl) View() string {
	config := config_module.GetConfig()

	var output string
	filteredChoices := m.FilteredChoices()

	questionIcon := lipgloss.NewStyle().Foreground(theme.Success).Render("?")
	checkIcon := lipgloss.NewStyle().Foreground(theme.Success).Render("✓")
	cursorIcon := lipgloss.NewStyle().Foreground(theme.Secondary).Bold(true).Render(">")
	paginationMessage := lipgloss.NewStyle().Foreground(theme.DarkenText).Render("\n(Move up and down to reveal more options | %d-%d of %d)")
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Secondary).Bold(true)
	borderStyle := lipgloss.NewStyle().
		Padding(1, 4).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(80)

	searchStyle := lipgloss.NewStyle().Foreground(theme.DarkenText)
	var searchMessage string
	if m.searchTerm != "" {
		searchMessage = searchStyle.Render(m.searchTerm)
	} else {
		searchMessage = searchStyle.Render("(Use arrow keys or type to search)")
	}

	var titleMessage string
	if m.Messages.Title != "" {
		titleMessage = questionIcon + " " + m.Messages.Title + ": " + searchMessage + "\n\n"
	} else {
		titleMessage = questionIcon + " Select an option:" + searchMessage + "\n\n"
	}

	if !config.HideLogomark {
		output += lipgloss.NewStyle().Foreground(theme.Primary).Render(logoascii.GetLogo(".body-builder")) + "\n"
	}

	if m.Border {
		output += titleMessage
	} else {
		output += strings.Replace(titleMessage, "\n\n", "\n", 1)
	}

	noOptionsStyle := lipgloss.NewStyle().Foreground(theme.Warn).Bold(true)
	if len(filteredChoices) == 0 {
		if m.Messages.EmptyOptions != "" {
			output += noOptionsStyle.Render(m.Messages.EmptyOptions + "\n")
		} else {
			output += noOptionsStyle.Render("No options found...\n")
		}
	} else {
		start := 0
		if m.cursor >= m.PerPage {
			start = m.cursor - m.PerPage + 1
		}
		end := min(start+m.PerPage, len(filteredChoices))

		for i := start; i < end; i++ {
			if m.cursor == i {
				output += fmt.Sprintf("%s %s\n", cursorIcon, selectedStyle.Render(filteredChoices[i]))
			} else {
				output += fmt.Sprintf("%s %s\n", " ", filteredChoices[i])
			}
		}
		if len(filteredChoices) > m.PerPage {
			output += fmt.Sprintf(paginationMessage, start+1, end, len(filteredChoices))
		}
	}

	if m.selected != -1 {
		var resultTitle string
		if m.Messages.Title != "" {
			resultTitle = m.Messages.Title
		} else {
			resultTitle = "Selected option:"
		}
		output = checkIcon + " " + resultTitle + ": " + selectedStyle.Bold(false).Render(filteredChoices[m.selected])
		return output + "\n"
	}

	if m.Border {
		return borderStyle.Render(output) + "\n"
	}
	return output + "\n"
}

func (m MenuImpl) Init() tea.Cmd {
	return nil
}

func (m MenuImpl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		filtered := m.FilteredChoices()
		switch keyMsg.Type {
		case tea.KeyRunes:
			m.searchTerm += keyMsg.String()
			m.cursor = 0
		case tea.KeyBackspace:
			if len(m.searchTerm) > 0 {
				m.searchTerm = m.searchTerm[:len(m.searchTerm)-1]
				m.cursor = 0
			}
		case tea.KeySpace:
			m.searchTerm += " "
			m.cursor = 0
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp, tea.KeyCtrlP:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown, tea.KeyCtrlN:
			if m.cursor < len(filtered)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			if len(filtered) > 0 {
				for i, c := range m.Choices {
					if c == filtered[m.cursor] {
						m.selected = i
						break
					}
				}
			}
			return m, tea.Quit
		}
	}
	return m, nil
}
