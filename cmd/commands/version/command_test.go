package version_command

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestExecutor(t *testing.T) {
	var exitCode int
	Exit = func(code int) {
		exitCode = code
	}

	rootCmd := &cobra.Command{
		Use: "test",
	}

	Executor(rootCmd)

	err := rootCmd.Execute()
	assert.NoError(t, err, "should no error when executing the root command")
	assert.Equal(t, 0, exitCode, "should exit with code 0")

	var found bool
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			found = true
			assert.Equal(t, "Show the version of the application", cmd.Short, "should usage description for 'version' command to be 'Show the version of the application'")
			break
		}
	}
	assert.True(t, found, "should 'version' command to be added")
}

func TestRunFunction(t *testing.T) {
	var exitCode int
	Exit = func(code int) {
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

	Executor(rootCmd)

	var command *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			command = cmd
			assert.Equal(t, "Show the version of the application", cmd.Short, "should usage description for 'version' command to be 'Show the version of the application'")
			break
		}
	}
	assert.True(t, command != nil, "should 'version' command to be found")

	command.Run(command, []string{})

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
