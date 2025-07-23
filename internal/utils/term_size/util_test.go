package term_size

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTerminalWidth_DefaultFunc(t *testing.T) {
	_ = GetTerminalWidth(100)
}

func TestGetTerminalWidth(t *testing.T) {
	origGetSizeFunc := getSizeFunc
	defer func() { getSizeFunc = origGetSizeFunc }()

	getSizeFunc = func() (int, error) { return 80, nil }
	width := GetTerminalWidth(100)
	if width != 80 {
		t.Errorf("expected 80, got %d", width)
	}

	getSizeFunc = func() (int, error) { return 120, nil }
	width = GetTerminalWidth(100)
	if width != 100 {
		t.Errorf("expected 100 to width > 100, got %d", width)
	}

	getSizeFunc = func() (int, error) { return 0, nil }
	width = GetTerminalWidth(100)
	if width != 100 {
		t.Errorf("expected 100 for width <= 0, got %d", width)
	}

	getSizeFunc = func() (int, error) { return 0, assert.AnError }
	width = GetTerminalWidth(100)
	if width != 100 {
		t.Errorf("expected 100 for error, got %d", width)
	}
}

func TestGetTerminalHeight_DefaultFunc(t *testing.T) {
	_ = GetTerminalHeight(100)
}

func TestGetTerminalHeight(t *testing.T) {
	origGetHeightFunc := GetHeightFunc
	defer func() { GetHeightFunc = origGetHeightFunc }()

	GetHeightFunc = func() (int, error) { return 50, nil }
	height := GetTerminalHeight(100)
	assert.Equal(t, 50, height)

	GetHeightFunc = func() (int, error) { return 120, nil }
	height = GetTerminalHeight(100)
	assert.Equal(t, 100, height)

	GetHeightFunc = func() (int, error) { return 0, nil }
	height = GetTerminalHeight(100)
	assert.Equal(t, 100, height)

	GetHeightFunc = func() (int, error) { return 0, assert.AnError }
	height = GetTerminalHeight(100)
	assert.Equal(t, 100, height)
}
