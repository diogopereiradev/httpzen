package help_command

import (
	"io"
	"os"
	"strings"
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

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w

	rootCmd.SetArgs([]string{"--help"})
	_ = rootCmd.Execute()

	w.Close()
	os.Stdout = origStdout
	outBytes, _ := io.ReadAll(r)
	output := string(outBytes)

	if !strings.Contains(output, "Httpzen CLI Tool for API Management and Development") {
		t.Errorf("Expected title in help output")
	}

	if !strings.Contains(output, "get") || !strings.Contains(output, "post") {
		t.Errorf("Expected commands in help output")
	}

	if strings.Contains(output, "completion") || strings.Contains(output, "help command") {
		t.Errorf("Should not show 'completion' or 'help' commands")
	}

	if !strings.Contains(output, "--testflag") {
		t.Errorf("Expected flag in help output")
	}

	if !strings.Contains(output, "Available parameters") {
		t.Errorf("Expected 'Available parameters' section")
	}
}

func TestInit_NoCommands(t *testing.T) {
	rootCmd := &cobra.Command{Use: "httpzen"}

	Init(rootCmd)

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w

	rootCmd.SetArgs([]string{"--help"})
	_ = rootCmd.Execute()

	w.Close()
	os.Stdout = origStdout
	outBytes, _ := io.ReadAll(r)
	output := string(outBytes)

	if !strings.Contains(output, "Usage") {
		t.Errorf("Expected usage in help output")
	}
}

func TestInit_NoFlags(t *testing.T) {
	rootCmd := &cobra.Command{Use: "httpzen"}
	cmd := &cobra.Command{Use: "get", Short: "Get resource"}
	rootCmd.AddCommand(cmd)

	Init(rootCmd)

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w

	rootCmd.SetArgs([]string{"--help"})
	_ = rootCmd.Execute()

	w.Close()
	os.Stdout = origStdout
	outBytes, _ := io.ReadAll(r)
	output := string(outBytes)

	if strings.Contains(output, "Available parameters") {
		if !strings.Contains(output, "--help") {
			t.Errorf("Section 'Available parameters' should only appear for default --help flag")
		}
	}
}
