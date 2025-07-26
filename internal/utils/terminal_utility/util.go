package terminal_utility

import (
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

var getSizeFunc = func() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	return width, err
}

var GetHeightFunc = func() (int, error) {
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	return height, err
}

var execCommand = func(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

var getGOOS = func() string {
	return runtime.GOOS
}

var clear map[string]func()

func GetTerminalWidth(max int) int {
	width, err := getSizeFunc()
	if err != nil || width <= 0 {
		return max
	}
	if width > max {
		return max
	}
	return width
}

func GetTerminalHeight(max int) int {
	height, err := GetHeightFunc()
	if err != nil || height <= 0 {
		return max
	}
	if height > max {
		return max
	}
	return height
}

func Clear() {
	value, ok := clear[getGOOS()]
	if ok {
		value()
	} else {
		return
	}
}

func init() {
	clear = make(map[string]func())

	clear["linux"] = linuxClear
	clear["windows"] = windowsClear
}

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
