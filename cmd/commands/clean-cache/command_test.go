package clean_cache_command

import (
	"testing"

	ip_cache_module "github.com/diogopereiradev/httpzen/internal/cache"
	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestInit_AddsCommand(t *testing.T) {
	rootCmd := &cobra.Command{Use: "root"}
	Init(rootCmd)
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "cleancache" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestCleancache_Run(t *testing.T) {
	var cleared, logged, exited bool
	oldClearCache := ip_cache_module.ClearCache
	oldSuccess := logger_module.Success
	oldExit := Exit
	defer func() {
		IpClearCache = oldClearCache
		LoggerSuccess = oldSuccess
		Exit = oldExit
	}()

	IpClearCache = func() { cleared = true }
	LoggerSuccess = func(msg string, width int) { logged = true }
	Exit = func(code int) { exited = code == 0 }

	rootCmd := &cobra.Command{Use: "root"}
	Init(rootCmd)

	var cleancacheCmd *cobra.Command

	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "cleancache" {
			cleancacheCmd = cmd
			break
		}
	}
	
	assert.NotNil(t, cleancacheCmd)
	cleancacheCmd.Run(cleancacheCmd, []string{})
	assert.True(t, cleared)
	assert.True(t, logged)
	assert.True(t, exited)
}
