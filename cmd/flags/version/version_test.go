package version_flag

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAddFlag(t *testing.T) {
	var exitCode int
	ExitFunc = func(code int) {
		exitCode = code
	}

	rootCmd := &cobra.Command{
		Use: "test",
	}

	AddFlag(rootCmd)

	err := rootCmd.Execute()
	assert.NoError(t, err, "should no error when executing the root command")
	assert.Equal(t, 0, exitCode, "should exit with code 0")

	flag := rootCmd.PersistentFlags().Lookup("version")

	assert.NotNil(t, flag, "should 'version' flag to be added")
	assert.Equal(t, "v", flag.Shorthand, "should shorthand for 'version' flag to be 'v'")
	assert.Equal(t, "Show the version of the application", flag.Usage, "should usage description for 'version' flag to be 'Show the version of the application'")
}

func TestRunFunction(t *testing.T) {
	var exitCode int
	ExitFunc = func(code int) {
		exitCode = code
	}

	Version = "1.0.0"
	BuildDate = "2025-07-19"
	Website = "https://example.com"
	Repository = "github.com/example/repo"
	License = "MIT"

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd := &cobra.Command{
		Use: "test",
	}

	AddFlag(rootCmd)

	rootCmd.Run(rootCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf []byte = make([]byte, 4096)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	assert.Contains(t, output, "Version: 1.0.0")
	assert.Contains(t, output, "Build Date: 2025-07-19")
	assert.Contains(t, output, "Website: https://example.com")
	assert.Contains(t, output, "Repository: github.com/example/repo")
	assert.Contains(t, output, "License: MIT")
	assert.Equal(t, 0, exitCode, "should exit with code 0")
}
