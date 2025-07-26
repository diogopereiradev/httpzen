package terminal_utility

import (
	"errors"
	"os/exec"
	"testing"
)

func TestGetTerminalWidth(t *testing.T) {
	oldGetWidthFunc := GetWidthFunc
	defer func() { GetWidthFunc = oldGetWidthFunc }()

	GetWidthFunc = func() (int, error) { return 80, nil }
	if w := GetTerminalWidth(100); w != 80 {
		t.Errorf("expected 80, got %d", w)
	}

	GetWidthFunc = func() (int, error) { return 120, nil }
	if w := GetTerminalWidth(100); w != 100 {
		t.Errorf("expected 100, got %d", w)
	}

	GetWidthFunc = func() (int, error) { return -1, errors.New("err") }
	if w := GetTerminalWidth(50); w != 50 {
		t.Errorf("expected 50, got %d", w)
	}
}

func TestGetTerminalHeight(t *testing.T) {
	oldGetHeightFunc := GetHeightFunc
	defer func() { GetHeightFunc = oldGetHeightFunc }()

	GetHeightFunc = func() (int, error) { return 24, nil }
	if h := GetTerminalHeight(30); h != 24 {
		t.Errorf("expected 24, got %d", h)
	}

	GetHeightFunc = func() (int, error) { return 40, nil }
	if h := GetTerminalHeight(30); h != 30 {
		t.Errorf("expected 30, got %d", h)
	}

	GetHeightFunc = func() (int, error) { return -1, errors.New("err") }
	if h := GetTerminalHeight(10); h != 10 {
		t.Errorf("expected 10, got %d", h)
	}
}

func TestClear_Linux(t *testing.T) {
	oldGetGOOS := getGOOS
	oldExecCommand := execCommand
	defer func() {
		getGOOS = oldGetGOOS
		execCommand = oldExecCommand
	}()

	getGOOS = func() string { return "linux" }
	called := false
	execCommand = func(name string, arg ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	clear["linux"] = func() { called = true }
	Clear()
	if !called {
		t.Error("expected linuxClear to be called")
	}
}

func TestLinuxClear(t *testing.T) {
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	called := false
	execCommand = func(name string, arg ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	linuxClear()
	if !called {
		t.Error("expected execCommand to be called in linuxClear")
	}

	execCommand = oldExecCommand
	linuxClear()
}

func TestWindowsClear(t *testing.T) {
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	called := false
	execCommand = func(name string, arg ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	windowsClear()
	if !called {
		t.Error("expected execCommand to be called in windowsClear")
	}

	execCommand = oldExecCommand
	windowsClear()
}

func TestRuntimeGOOSExecution(t *testing.T) {
	_ = getGOOS()
}

func TestGetWidthFuncExecution(t *testing.T) {
	_, _ = GetWidthFunc()
}

func TestGetHeightFuncExecution(t *testing.T) {
	_, _ = GetHeightFunc()
}