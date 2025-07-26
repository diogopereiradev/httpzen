package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestThemeColors(t *testing.T) {
	assert.Equal(t, lipgloss.AdaptiveColor{Light: "3", Dark: "3"}, Primary)
	assert.Equal(t, lipgloss.AdaptiveColor{Light: "6", Dark: "6"}, Secondary)
	assert.Equal(t, lipgloss.AdaptiveColor{Light: "7", Dark: "7"}, LightText)
	assert.Equal(t, lipgloss.AdaptiveColor{Light: "8", Dark: "8"}, DarkenText)
	assert.Equal(t, lipgloss.AdaptiveColor{Light: "1", Dark: "1"}, Error)
	assert.Equal(t, lipgloss.AdaptiveColor{Light: "11", Dark: "11"}, Warn)
	assert.Equal(t, lipgloss.AdaptiveColor{Light: "2", Dark: "2"}, Success)
}
