package blink_cursor_component

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type BlinkCursor struct {
	char       rune
	visible     bool
	blinking bool
	blinkRate  time.Duration
	blinkTimer tea.Cmd
}

type BlinkCursorMsg struct {}

func initialModel() *BlinkCursor {
	return &BlinkCursor{
		char:       '|',
		visible:     true,
		blinking:    true,
		blinkRate:   500 * time.Millisecond,
		blinkTimer:  nil,
	}
}

func New() *BlinkCursor {
	m := initialModel()
	return m
}

func (m *BlinkCursor) Render() rune {
	cursorChar := ' '
	if m.blinking && m.visible {
		cursorChar = m.char
	}
	return cursorChar
}

func (m *BlinkCursor) Show() tea.Cmd {
	m.visible = true
	m.blinkTimer = tea.Tick(m.blinkRate, func(t time.Time) tea.Msg {
		if !m.visible {
			m.blinking = false
			return nil
		}
		m.blinking = !m.blinking
		return BlinkCursorMsg{}
	})
	return m.blinkTimer
}

func (m *BlinkCursor) Hide() {
	m.visible = false
	m.blinking = false
}