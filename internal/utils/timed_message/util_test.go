package timed_message_util

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewTimedMessage(t *testing.T) {
	tm := NewTimedMessage()
	assert.NotNil(t, tm)
	assert.Equal(t, "", tm.Message)
	assert.False(t, tm.Visible)
	assert.Equal(t, time.Duration(0), tm.Duration)
}

func TestShowAndRender(t *testing.T) {
	tm := NewTimedMessage()
	msg := "Hello, World!"
	dur := 10 * time.Millisecond
	cmd := tm.Show(msg, dur)

	assert.Equal(t, msg, tm.Message)
	assert.True(t, tm.Visible)
	assert.Equal(t, dur, tm.Duration)

	rendered := tm.Render()
	assert.Contains(t, rendered, msg)

	ch := make(chan tea.Msg, 1)
	go func() {
		ch <- cmd()
	}()
	select {
	case m := <-ch:
		_, ok := m.(TimedMessageExpiredMsg)
		assert.True(t, ok)
		assert.False(t, tm.Visible)

		assert.Equal(t, "", tm.Render())

	case <-time.After(100 * time.Millisecond):
		t.Fatal("TimedMessage did not expire in time")
	}
}

func TestRenderWhenNotVisible(t *testing.T) {
	tm := NewTimedMessage()
	tm.Message = "Should not show"
	tm.Visible = false
	assert.Equal(t, "", tm.Render())
}
