package timed_message_component

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

type TimedMessage struct {
	Message  string
	Visible  bool
	Duration time.Duration
}

type TimedMessageExpiredMsg struct{}

func NewTimedMessage() *TimedMessage {
	return &TimedMessage{}
}

func (d *TimedMessage) Show(message string, duration time.Duration) tea.Cmd {
	d.Message = message
	d.Visible = true
	d.Duration = duration
	return func() tea.Msg {
		time.Sleep(d.Duration)
		d.Visible = false
		return TimedMessageExpiredMsg{}
	}
}

func (d *TimedMessage) Render() string {
	if !d.Visible {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(theme.LightText).
		Background(theme.Success).
		Padding(0, 2).
		Bold(true).
		Margin(1, 0)

	return style.Render(d.Message)
}
