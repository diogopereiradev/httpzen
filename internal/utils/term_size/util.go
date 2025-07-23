package term_size

import (
	"os"

	"golang.org/x/term"
)

var getSizeFunc = func() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	return width, err
}

var getHeightFunc = func() (int, error) {
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	return height, err
}

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
	height, err := getHeightFunc()
	if err != nil || height <= 0 {
		return max
	}
	if height > max {
		return max
	}
	return height
}