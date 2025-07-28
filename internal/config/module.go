package config_module

import (
	"os"

	app_path_util "github.com/diogopereiradev/httpzen/internal/utils/app_path"
	"github.com/spf13/viper"
)

type Config struct {
	SlowResponseThreshold int  `json:"slow_response_threshold"`
	HideLogomark          bool `json:"hide_logomark"`
}

var CONFIG_NAME string = "config"
var CONFIG_EXTENSION string = "json"

var mkdirAll = os.MkdirAll

func GetConfig() Config {
	configPath := app_path_util.GetConfigPath()
	v := viper.New()
	v.SetConfigName(CONFIG_NAME)
	v.SetConfigType(CONFIG_EXTENSION)
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return InitConfig()
	}
	return Config{
		SlowResponseThreshold: v.GetInt("slow_response_threshold"),
		HideLogomark:          v.GetBool("hide_logomark"),
	}
}

func UpdateConfig(newConfig Config) error {
	v := viper.New()
	v.Set("slow_response_threshold", newConfig.SlowResponseThreshold)
	v.Set("hide_logomark", newConfig.HideLogomark)

	configPath := app_path_util.GetConfigPath()
	if err := mkdirAll(configPath, 0755); err != nil {
		return err
	}

	v.SetConfigName(CONFIG_NAME)
	v.SetConfigType(CONFIG_EXTENSION)
	v.AddConfigPath(configPath)

	if err := v.WriteConfigAs(configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION); err != nil {
		return err
	}
	return nil
}

func InitConfig() Config {
	config := Config{
		SlowResponseThreshold: 500,
		HideLogomark:          false,
	}

	configPath := app_path_util.GetConfigPath()
	if err := mkdirAll(configPath, 0755); err != nil {
		panic(err)
	}

	v := viper.New()
	v.SetDefault("slow_response_threshold", config.SlowResponseThreshold)
	v.SetDefault("hide_logomark", config.HideLogomark)

	v.SetConfigName(CONFIG_NAME)
	v.SetConfigType(CONFIG_EXTENSION)
	v.AddConfigPath(configPath)

	if err := v.WriteConfigAs(configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION); err != nil {
		panic(err)
	}
	return config
}

func GetConfigFilePath() string {
	configPath := app_path_util.GetConfigPath()
	return configPath + "/" + CONFIG_NAME + "." + CONFIG_EXTENSION
}
