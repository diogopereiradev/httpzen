package help_command

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var CategorizedFlags = map[string][]string{
	"Main parameters": {"help"},
	"Data":            {"header", "body"},
}

var CategorizedFlagsOrder = []string{
	"Main parameters",
	"Data",
}

func padRight(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(" ", length-len(str))
}

func renderCommands(cmds []*cobra.Command, maxUseLen int) string {
	var b strings.Builder
	for _, cmd := range cmds {
		if cmd.Use == "completion" || strings.HasPrefix(cmd.Use, "help") {
			continue
		}
		b.WriteString("  ")
		b.WriteString(padRight(cmd.Use, maxUseLen+4))
		b.WriteString(cmd.Short)
		b.WriteString("\n")
	}
	return b.String()
}

func renderFlags(cmd *cobra.Command, maxFlagNameLen int) string {
	var b strings.Builder
	shown := map[string]bool{}

	for i, category := range CategorizedFlagsOrder {
		flagNames := CategorizedFlags[category]
		hasAny := false
		for _, flagName := range flagNames {
			if cmd.Flags().Lookup(flagName) != nil {
				hasAny = true
				break
			}
		}

		if !hasAny {
			continue
		}

		b.WriteString(category + "\n\n")
		for _, flagName := range flagNames {
			flag := cmd.Flags().Lookup(flagName)
			shown[flagName] = true

			pad := padRight("", maxFlagNameLen-len(flagName)+4)
			if flag.Shorthand != "" {
				b.WriteString("  --" + flag.Name + ", -" + flag.Shorthand + pad + flag.Usage + "\n")
			} else {
				b.WriteString("  --" + flag.Name + "    " + pad + flag.Usage + "\n")
			}
		}
		if i < len(CategorizedFlagsOrder)-1 {
			b.WriteString("\n")
		}
	}

	var uncategorized []string
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !shown[f.Name] {
			uncategorized = append(uncategorized, f.Name)
		}
	})

	if len(uncategorized) > 0 {
		b.WriteString("Available parameters\n\n")
		for _, flagName := range uncategorized {
			flag := cmd.Flags().Lookup(flagName)

			pad := padRight("", maxFlagNameLen-len(flagName)+4)
			if flag.Shorthand != "" {
				b.WriteString("  --" + flag.Name + ", -" + flag.Shorthand + pad + flag.Usage + "\n")
			} else {
				b.WriteString("  --" + flag.Name + "    " + pad + flag.Usage + "\n")
			}
		}
	}
	return b.String()
}

func Init(rootCmd *cobra.Command) {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		borderStyle := lipgloss.
			NewStyle().
			Width(term_size.GetTerminalWidth(90)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(0, 3)

		logoStyle := lipgloss.
			NewStyle().
			Foreground(theme.Primary).
			Bold(true)

		titleStyle := lipgloss.
			NewStyle().
			Foreground(theme.Secondary).
			Bold(true).
			Underline(true)

		fieldStyle := lipgloss.
			NewStyle().
			Foreground(theme.DarkenText).
			Bold(true)

		var b strings.Builder
		b.WriteString(logoStyle.Render(logoascii.GetLogo(".help")) + "\n")
		b.WriteString(titleStyle.Render("Httpzen CLI Tool for API Management and Development") + "\n\n")
		b.WriteString("Usage\n")
		b.WriteString("  $ httpzen " + fieldStyle.Render("[METHOD] ") + fieldStyle.Render("[URL] ") + fieldStyle.Render("[PARAMETERS ...]"))

		cmds := cmd.Commands()
		if len(cmds) > 0 {
			b.WriteString("\n\nMain Commands\n\n")
		}

		maxUseLen := 0
		for _, c := range cmds {
			if c.Use == "completion" || strings.HasPrefix(c.Use, "help") {
				continue
			}
			if len(c.Use) > maxUseLen {
				maxUseLen = len(c.Use)
			}
		}
		b.WriteString(renderCommands(cmds, maxUseLen))
		b.WriteString("\n")

		maxFlagNameLen := 0
		for _, category := range CategorizedFlagsOrder {
			flagNames := CategorizedFlags[category]
			for _, flagName := range flagNames {
				if len(flagName) > maxFlagNameLen {
					maxFlagNameLen = len(flagName)
				}
			}
		}
		b.WriteString(renderFlags(cmd, maxFlagNameLen))

		fmt.Println(borderStyle.Render(b.String()))
	})
}
