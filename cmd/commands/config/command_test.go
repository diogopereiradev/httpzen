package config_command

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestInit_AddsConfigCommand(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	Init(rootCmd)
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "config" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Config command was not added to root command")
	}
}

func TestConfigCommand_Run_Error(t *testing.T) {
	var called bool
	oldRunConfigEditor := runConfigEditor
	runConfigEditor = func() error {
		called = true
		return errors.New("fail")
	}
	defer func() { runConfigEditor = oldRunConfigEditor }()

	rootCmd := &cobra.Command{Use: "root"}
	Init(rootCmd)
	cmd, _, err := rootCmd.Find([]string{"config"})
	if err != nil {
		t.Fatalf("Could not find config command: %v", err)
	}
	cmd.Run(cmd, []string{})
	if !called {
		t.Errorf("runConfigEditor was not called")
	}
}

func TestConfigCommand_Run_Success(t *testing.T) {
	var called bool
	oldRunConfigEditor := runConfigEditor
	runConfigEditor = func() error {
		called = true
		return nil
	}
	defer func() { runConfigEditor = oldRunConfigEditor }()

	rootCmd := &cobra.Command{Use: "root"}
	Init(rootCmd)
	cmd, _, err := rootCmd.Find([]string{"config"})
	if err != nil {
		t.Fatalf("Could not find config command: %v", err)
	}
	cmd.Run(cmd, []string{})
	if !called {
		t.Errorf("runConfigEditor was not called")
	}
}
