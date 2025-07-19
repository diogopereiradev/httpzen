package app_config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetConfigPath
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
		if key == "APPDATA" { return "C:\\Users\\TestUser\\AppData\\Roaming" }
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

// GetConfig
func TestGetConfig(t *testing.T) {
	CONFIG_NAME = "httpzen_test"

	configPath := GetConfigPath()
	configFile := configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION
	_ = os.Remove(configFile)
	_, err := os.Stat(configFile)

	assert.Error(t, err, "should config file not exist before test")

	config := GetConfig()
	assert.Equal(t, 500, config.SlowResponseThreshold, "should have SlowResponseThreshold set to default value")

	removeErr := os.Remove(configFile)
	assert.NoError(t, removeErr, "should not fail to remove test config file")
}

func TestGetConfigReadInConfig(t *testing.T) {
	CONFIG_NAME = "httpzen_test"

	configPath := GetConfigPath()
	configFile := configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION
	InitConfig()

	config := GetConfig()
	assert.Equal(t, 500, config.SlowResponseThreshold, "should have SlowResponseThreshold set to default value")

	removeErr := os.Remove(configFile)
	assert.NoError(t, removeErr, "should not fail to remove test config file")
}

// InitConfig
func TestInitConfig(t *testing.T) {
	CONFIG_NAME = "httpzen_test"

	configPath := GetConfigPath()
	configFile := configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION
	_ = os.Remove(configFile)

	config := InitConfig()

	_, err := os.Stat(configFile)
	assert.NoError(t, err, "should config file be created by InitConfig")
	assert.Equal(t, 500, config.SlowResponseThreshold, "should have default SlowResponseThreshold")

	removeErr := os.Remove(configFile)
	assert.NoError(t, removeErr, "should not fail to remove test config file")
}

func TestInitConfigPanicOnMkdirAll(t *testing.T) {
	originalMkdirAll := mkdirAll
	mkdirAll = func (string, os.FileMode) error { return assert.AnError }
	defer func() { mkdirAll = originalMkdirAll }()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not occur")
		}
	}()
	InitConfig()
}

func TestInitConfigPanicOnWriteConfigAs(t *testing.T) {
	configPathBkp := CONFIG_NAME
	CONFIG_NAME = "httpzen_readonly_test"

	configPath := GetConfigPath()
	_ = os.MkdirAll(configPath, 0755)

	configFile := configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION
	_ = os.WriteFile(configFile, []byte("{}"), 0444)

	err := os.Chmod(configFile, 0444)
	if err != nil {
		t.Skip("Not is possible to set permissions for test")
	}

	defer func() {
		CONFIG_NAME = configPathBkp
		os.Chmod(configFile, 0644)
		os.RemoveAll(configPath)

		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not occur")
		}
	}()

	InitConfig()
}

// UpdateConfig
func TestUpdateConfig(t *testing.T) {
	CONFIG_NAME = "httpzen_test"

	configPath := GetConfigPath()
	configFile := configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION
	_ = os.Remove(configFile)

	InitConfig()
	configData := Config{
		SlowResponseThreshold: 1000,
	}
	err := UpdateConfig(configData)
	assert.NoError(t, err, "should not fail to update config")

	updatedConfig := GetConfig()
	assert.Equal(t, 1000, updatedConfig.SlowResponseThreshold, "should have SlowResponseThreshold updated")

	removeErr := os.Remove(configFile)
	assert.NoError(t, removeErr, "should not fail to remove test config file")
}

func TestUpdateConfigPanicOnMkdirAll(t *testing.T) {
	originalMkdirAll := mkdirAll
	mkdirAll = func (string, os.FileMode) error { return assert.AnError }
	defer func() { mkdirAll = originalMkdirAll }()

	config := Config{
		SlowResponseThreshold: 1000,
	}
	UpdateConfig(config)
}

func TestUpdateConfigPanicOnWriteConfigAs(t *testing.T) {
	configPathBkp := CONFIG_NAME
	CONFIG_NAME = "httpzen_readonly_update_test"

	configPath := GetConfigPath()
	_ = os.MkdirAll(configPath, 0755)

	configFile := configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION
	_ = os.WriteFile(configFile, []byte("{}"), 0444)

	err := os.Chmod(configFile, 0444)
	if err != nil {
		t.Skip("Not is possible to set permissions for test")
	}

	defer func() {
		CONFIG_NAME = configPathBkp
		os.Chmod(configFile, 0644)
		os.RemoveAll(configPath)
	}()

	config := Config{
		SlowResponseThreshold: 1000,
	}
	UpdateConfig(config)
}
