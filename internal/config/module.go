package app_config

import (
	"os"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	SlowResponseThreshold int `json:"slow_response_threshold"`
}

var CONFIG_NAME string = "httpzen"
var CONFIG_EXTENSION string = "json"
var GOOS = runtime.GOOS

var userHomeDir = os.UserHomeDir
var getenv = os.Getenv
var mkdirAll = os.MkdirAll

func GetConfig() Config {
	configPath := GetConfigPath()

	viper.SetConfigName(CONFIG_NAME)
	viper.SetConfigType(CONFIG_EXTENSION)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return InitConfig()
	}
	return Config{
		SlowResponseThreshold: viper.GetInt("slow_response_threshold"),
	}
}

func UpdateConfig(newConfig Config) error {
	viper.Set("slow_response_threshold", newConfig.SlowResponseThreshold)

	configPath := GetConfigPath()
	if err := mkdirAll(configPath, 0755); err != nil {
		return err
	}

	if err := viper.WriteConfigAs(configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION); err != nil {
		return err
	}
	return nil
}

func GetConfigPath() string {
	home, err := userHomeDir()
	if err == nil {
		switch GOOS {
		case "windows":
			appData := getenv("APPDATA")
			if appData != "" {
				return appData + "\\httpzen"
			}
		case "darwin":
			return home + "/Library/Application Support/httpzen"
		default:
			return home + "/.config/httpzen"
		}
	}
	panic(err)
}

func InitConfig() Config {
	config := Config{
		SlowResponseThreshold: 500,
	}

	configPath := GetConfigPath()
	if err := mkdirAll(configPath, 0755); err != nil {
		panic(err)
	}

	viper.SetDefault("slow_response_threshold", config.SlowResponseThreshold)

	if err := viper.WriteConfigAs(configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION); err != nil {
		panic(err)
	}
	return config
}
