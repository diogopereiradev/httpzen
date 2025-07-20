package app_path_util

import (
	"os"
	"runtime"
)

var GOOS = runtime.GOOS

var getenv = os.Getenv
var userHomeDir = os.UserHomeDir

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
