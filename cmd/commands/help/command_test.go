package help_command

import (
	"bytes"
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPadRight(t *testing.T) {
	assert.Equal(t, "foo   ", padRight("foo", 6))
	assert.Equal(t, "foobar", padRight("foobar", 6))
	assert.Equal(t, "foobar", padRight("foobar", 3))
}

func TestRenderCommands(t *testing.T) {
	cmd1 := &cobra.Command{Use: "get", Short: "Get resource"}
	cmd2 := &cobra.Command{Use: "post", Short: "Post resource"}
	cmd3 := &cobra.Command{Use: "completion", Short: "Should be skipped"}
	cmd4 := &cobra.Command{Use: "help", Short: "Should be skipped"}
	cmds := []*cobra.Command{cmd1, cmd2, cmd3, cmd4}
	out := renderCommands(cmds, 4)

	expected := "  get     Get resource\n  post    Post resource\n"
	assert.Equal(t, expected, out)
}

func TestRenderFlags_CategorizedAndUncategorized(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("help", "", "Show help")
	cmd.Flags().String("header", "", "Set header")
	cmd.Flags().String("body", "", "Set body")
	cmd.Flags().String("extra", "", "Extra flag")
	cmd.Flags().StringP("short", "s", "", "Short flag")

	out := renderFlags(cmd, 6)
	assert.Contains(t, out, "Main parameters")
	assert.Contains(t, out, "Data")
	assert.Contains(t, out, "Available parameters")
	assert.Contains(t, out, "--extra")
	assert.Contains(t, out, "--short, -s")
}

func TestRenderFlags_OnlyCategorized(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("help", "", "Show help")
	cmd.Flags().String("header", "", "Set header")
	cmd.Flags().String("body", "", "Set body")

	out := renderFlags(cmd, 6)
	assert.Contains(t, out, "Main parameters")
	assert.Contains(t, out, "Data")
	assert.NotContains(t, out, "Available parameters")
}

func TestRenderFlags_OnlyUncategorized(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("foo", "", "Foo flag")
	cmd.Flags().String("bar", "", "Bar flag")

	out := renderFlags(cmd, 6)
	assert.Contains(t, out, "Available parameters")
	assert.Contains(t, out, "--foo")
	assert.Contains(t, out, "--bar")
}

func TestRenderFlags_NoFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}

	out := renderFlags(cmd, 6)
	assert.Equal(t, "", out)
}

func TestInit_SetsHelpFunc(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringP("withshort", "w", "", "With shorthand")
	cmd.Flags().String("noshort", "", "No shorthand")

	CategorizedFlags["Main parameters"] = []string{"withshort", "noshort"}
	CategorizedFlagsOrder = []string{"Main parameters"}

	out := renderFlags(cmd, 8)
	assert.Contains(t, out, "--withshort, -w   With shorthand")
	assert.Contains(t, out, "--noshort         No shorthand")
}

func TestInit_Integration(t *testing.T) {
	oldGetTerminalWidth := terminal_utility.GetWidthFunc
	terminal_utility.GetWidthFunc = func() (int, error) { return 40, nil }
	defer func() { terminal_utility.GetWidthFunc = oldGetTerminalWidth }()

	oldPrimary := theme.Primary
	oldSecondary := theme.Secondary
	oldDarkenText := theme.DarkenText
	theme.Primary = lipgloss.AdaptiveColor{Light: "1", Dark: "1"}
	theme.Secondary = lipgloss.AdaptiveColor{Light: "2", Dark: "2"}
	theme.DarkenText = lipgloss.AdaptiveColor{Light: "3", Dark: "3"}
	defer func() {
		theme.Primary = oldPrimary
		theme.Secondary = oldSecondary
		theme.DarkenText = oldDarkenText
	}()

	oldCategorizedFlags := make(map[string][]string)
	for k, v := range CategorizedFlags {
		copied := make([]string, len(v))
		copy(copied, v)
		oldCategorizedFlags[k] = copied
	}
	oldCategorizedFlagsOrder := make([]string, len(CategorizedFlagsOrder))
	copy(oldCategorizedFlagsOrder, CategorizedFlagsOrder)
	defer func() {
		CategorizedFlags = oldCategorizedFlags
		CategorizedFlagsOrder = oldCategorizedFlagsOrder
	}()

	cmd := &cobra.Command{Use: "root", Short: "Root command"}
	child := &cobra.Command{Use: "child", Short: "Child command"}

	cmd.AddCommand(child)
	cmd.Flags().String("help", "", "Show help")
	cmd.Flags().String("header", "", "Set header")
	cmd.Flags().String("body", "", "Set body")
	cmd.Flags().String("extra", "", "Extra flag")
	cmd.Flags().StringP("short", "s", "", "Short flag")

	CategorizedFlags = map[string][]string{
		"Main parameters": {"help"},
		"Data":            {"header", "body"},
	}
	CategorizedFlagsOrder = []string{"Main parameters", "Data"}

	Init(cmd)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	cmd.Help()
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	assert.Contains(t, output, "/_/                     │")
	assert.Contains(t, output, "Httpzen CLI Tool for API")
	assert.Contains(t, output, "Main Commands")
	assert.Contains(t, output, "child    Child command")
	assert.Contains(t, output, "Main parameters")
	assert.Contains(t, output, "Data")
	assert.Contains(t, output, "Available parameters")
	assert.Contains(t, output, "--extra")
	assert.Contains(t, output, "--short, -s")
}

func TestInit_CommandsWithCompletionAndHelpAreSkipped(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.AddCommand(&cobra.Command{Use: "completion", Short: "Should skip"})
	cmd.AddCommand(&cobra.Command{Use: "helpSomething", Short: "Should skip"})
	cmd.AddCommand(&cobra.Command{Use: "normal", Short: "Normal"})

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Init(cmd)

	cmd.Help()
	w.Close()
	os.Stdout = oldStdout

	_, _ = buf.ReadFrom(r)
	output := buf.String()

	assert.Contains(t, output, "normal    Normal")
	assert.NotContains(t, output, "completion")
	assert.NotContains(t, output, "helpSomething")
}
