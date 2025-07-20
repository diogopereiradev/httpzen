package app_path_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigPath(t *testing.T) {
	GOOS = "linux"

	originalUserHomeDir := userHomeDir
	originalGetenv := getenv
	userHomeDir = func() (string, error) { return "/home/testuser", nil }
	getenv = func(key string) string { return "" }
	defer func() {
		userHomeDir = originalUserHomeDir
		getenv = originalGetenv
	}()

	configPath := GetConfigPath()
	assert.NotEmpty(t, configPath, "should return a non-empty config path")
}

func TestGetConfigPathWindows(t *testing.T) {
	GOOS = "windows"

	originalUserHomeDir := userHomeDir
	originalGetenv := getenv
	userHomeDir = func() (string, error) { return "C:\\Users\\TestUser", nil }
	getenv = func(key string) string {
		if key == "APPDATA" {
			return "C:\\Users\\TestUser\\AppData\\Roaming"
		}
		return ""
	}
	defer func() {
		userHomeDir = originalUserHomeDir
		getenv = originalGetenv
	}()

	configPath := GetConfigPath()
	assert.Equal(t, "C:\\Users\\TestUser\\AppData\\Roaming\\httpzen", configPath, "should return correct path for Windows")
}

func TestGetConfigPathDarwin(t *testing.T) {
	GOOS = "darwin"

	originalUserHomeDir := userHomeDir
	originalGetenv := getenv
	userHomeDir = func() (string, error) { return "/Users/TestUser", nil }
	getenv = func(key string) string { return "" }
	defer func() {
		userHomeDir = originalUserHomeDir
		getenv = originalGetenv
	}()

	configPath := GetConfigPath()
	assert.Equal(t, "/Users/TestUser/Library/Application Support/httpzen", configPath, "should return correct path for macOS")
}

func TestGetConfigPathLinux(t *testing.T) {
	GOOS = "linux"

	originalUserHomeDir := userHomeDir
	originalGetenv := getenv

	userHomeDir = func() (string, error) { return "/home/testuser", nil }
	getenv = func(key string) string { return "" }
	defer func() {
		userHomeDir = originalUserHomeDir
		getenv = originalGetenv
	}()

	configPath := GetConfigPath()
	assert.Equal(t, "/home/testuser/.config/httpzen", configPath, "should return correct path for Linux")
}

func TestGetConfigPathError(t *testing.T) {
	originalUserHomeDir := userHomeDir
	originalGetenv := getenv
	userHomeDir = func() (string, error) { return "", assert.AnError }
	getenv = func(_ string) string { return "" }
	defer func() {
		userHomeDir = originalUserHomeDir
		getenv = originalGetenv
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not occur")
		}
	}()
	GetConfigPath()
}
