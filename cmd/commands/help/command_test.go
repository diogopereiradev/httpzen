package help_command

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestInit_HelpOutput(t *testing.T) {
	rootCmd := &cobra.Command{Use: "httpzen"}
	cmd1 := &cobra.Command{Use: "get", Short: "Get resource"}
	cmd2 := &cobra.Command{Use: "post", Short: "Post resource"}
	cmd3 := &cobra.Command{Use: "completion", Short: "Shell completion"}
	cmd4 := &cobra.Command{Use: "help", Short: "Help command"}

	rootCmd.AddCommand(cmd1, cmd2, cmd3, cmd4)
	rootCmd.Flags().Bool("testflag", false, "A test flag")

	Init(rootCmd)

	rootCmd.SetArgs([]string{"--help"})
	_ = rootCmd.Execute()
}

func Test_renderFlags_categorized(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Bool("json", false, "Test JSON flag")
	cmd.Flags().Bool("raw", false, "Test raw flag")
	cmd.Flags().Bool("form", false, "Test form flag")
	cmd.Flags().Bool("multipart", false, "Test multipart flag")
	cmd.Flags().Bool("headers", false, "Test headers flag")
	cmd.Flags().Bool("body", false, "Test body flag")
	cmd.Flags().Bool("meta", false, "Test meta flag")

	maxFlagNameLen := 10
	result := renderFlags(cmd, maxFlagNameLen)

	if result == "" {
		t.Errorf("Expected no output, got:\n%s", result)
	}
}

func Test_renderFlags_uncategorized(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Bool("uncategorized-foobar", false, "Test uncategorized flag")
	cmd.Flags().BoolP("uncategorized-foobarp", "Z", false, "Test uncategorized flag")

	maxFlagNameLen := 10
	result := renderFlags(cmd, maxFlagNameLen)

	if result == "" {
		t.Errorf("Expected no output, got:\n%s", result)
	}
}

func Test_padRight(t *testing.T) {
	result := padRight("test", 10)
	expected := "test      "
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	result = padRight("longertext", 10)
	expected = "longertext"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
