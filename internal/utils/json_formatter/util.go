package json_formatter

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func FormatJSON(input string) string {
	var prettyBuf = new(bytes.Buffer)
	if err := json.Indent(prettyBuf, []byte(input), "", "  "); err != nil {
		return input
	}
	jsonStr := strings.TrimSpace(prettyBuf.String())

	var result strings.Builder
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))    // Yellow
	stringStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Green
	numberStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // Light Blue
	boolStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))   // Magenta
	nullStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))   // Red

	i := 0
	for i < len(jsonStr) {
		c := jsonStr[i]
		if c == '"' {
			// Key or string value detection
			j := i + 1
			for j < len(jsonStr) && jsonStr[j] != '"' {
				if jsonStr[j] == '\\' && j+1 < len(jsonStr) {
					j++ // Skip escape character
				}
				j++
			}
			if j < len(jsonStr) {
				segment := jsonStr[i : j+1]
				// Check if it's a key (comes before ':')
				k := j + 1
				for k < len(jsonStr) && (jsonStr[k] == ' ' || jsonStr[k] == '\n' || jsonStr[k] == '\t') {
					k++
				}
				if k < len(jsonStr) && jsonStr[k] == ':' {
					result.WriteString(keyStyle.Render(segment))
				} else {
					result.WriteString(stringStyle.Render(segment))
				}
				i = j + 1
				continue
			}
		}
		// Numbers
		if (c >= '0' && c <= '9') || c == '-' {
			j := i
			for j < len(jsonStr) && ((jsonStr[j] >= '0' && jsonStr[j] <= '9') || jsonStr[j] == '.' || jsonStr[j] == 'e' || jsonStr[j] == 'E' || jsonStr[j] == '-') {
				j++
			}
			result.WriteString(numberStyle.Render(jsonStr[i:j]))
			i = j
			continue
		}
		// Booleans and null
		if strings.HasPrefix(jsonStr[i:], "true") {
			result.WriteString(boolStyle.Render("true"))
			i += 4
			continue
		}
		if strings.HasPrefix(jsonStr[i:], "false") {
			result.WriteString(boolStyle.Render("false"))
			i += 5
			continue
		}
		if strings.HasPrefix(jsonStr[i:], "null") {
			result.WriteString(nullStyle.Render("null"))
			i += 4
			continue
		}
		// Other characters
		result.WriteByte(c)
		i++
	}
	return result.String()
}
