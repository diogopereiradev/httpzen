package html_formatter

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yosssi/gohtml"
)

func FormatHTML(input string) string {
	formatted := gohtml.Format(input)
	lines := strings.Split(formatted, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		colored := trimmed
		colored = colorTags(colored)
		colored = colorAttrs(colored)
		colored = colorVals(colored)
		result = append(result, colored)
	}
	return strings.Join(result, "\n")
}

func colorTags(s string) string {
	tagRegex := regexp.MustCompile(`(<\/?[a-zA-Z0-9]+)`)
	return tagRegex.ReplaceAllStringFunc(s, func(tag string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(tag)
	})
}

func colorAttrs(s string) string {
	return colorAttrsWithRegex(s, regexp.MustCompile(`([a-zA-Z\-]+)(=)`))
}

func colorAttrsWithRegex(s string, attrRegex *regexp.Regexp) string {
	return attrRegex.ReplaceAllStringFunc(s, func(attr string) string {
		parts := attrRegex.FindStringSubmatch(attr)
		if len(parts) == 3 {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render(parts[1]) + parts[2]
		}
		return attr
	})
}

func colorVals(s string) string {
	valRegex := regexp.MustCompile(`("[^"]*")`)
	return valRegex.ReplaceAllStringFunc(s, func(val string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(val)
	})
}
