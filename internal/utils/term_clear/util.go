package term_clear

import (
	"os"
	"os/exec"
	"runtime"
)

var execCommand = func(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

var getGOOS = func() string {
	return runtime.GOOS
}

var clear map[string]func()

func linuxClear() {
	cmd := execCommand("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func windowsClear() {
	cmd := execCommand("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func init() {
	clear = make(map[string]func())

	clear["linux"] = linuxClear
	clear["windows"] = windowsClear
}

func Clear() {
	value, ok := clear[getGOOS()]
	if ok {
		value()
	} else {
		return
	}
}
