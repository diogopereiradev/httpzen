package theme

import "github.com/charmbracelet/lipgloss"

var (
	Primary = lipgloss.AdaptiveColor{ Light: "3", Dark: "3" }
	Secondary = lipgloss.AdaptiveColor{ Light: "6", Dark: "6" }

	LightText = lipgloss.AdaptiveColor{ Light: "7", Dark: "7" }
	DarkenText = lipgloss.AdaptiveColor{ Light: "8", Dark: "8" }

	Error = lipgloss.AdaptiveColor{ Light: "1", Dark: "1" }
	Warn = lipgloss.AdaptiveColor{ Light: "11", Dark: "11" }
	Success = lipgloss.AdaptiveColor{ Light: "2", Dark: "2" }
)