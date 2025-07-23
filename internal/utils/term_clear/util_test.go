package term_clear

import (
	"os/exec"
	"runtime"
	"testing"
)

func TestExecCommand(t *testing.T) {
	_ = execCommand("test", "arg1", "arg2")
}

func TestClear_Linux(t *testing.T) {
	called := false
	execCommand = func(name string, arg ...string) *exec.Cmd {
		called = true
		if name != "clear" {
			t.Errorf("expected 'clear', got '%s'", name)
		}
		return exec.Command("echo")
	}
	runtimeGOOS := runtime.GOOS
	clear["linux"]()
	if !called && runtimeGOOS == "linux" {
		t.Error("linuxClear was not called")
	}
}

func TestClear_Windows(t *testing.T) {
	called := false
	execCommand = func(name string, arg ...string) *exec.Cmd {
		called = true
		if name != "cmd" || len(arg) < 2 || arg[0] != "/c" || arg[1] != "cls" {
			t.Errorf("expected 'cmd /c cls', got '%s %v'", name, arg)
		}
		return exec.Command("echo")
	}
	runtimeGOOS := runtime.GOOS
	clear["windows"]()
	if !called && runtimeGOOS == "windows" {
		t.Error("windowsClear was not called")
	}
}

func TestClear_Unknown(t *testing.T) {
	called := false
	clear["unknown"] = func() { called = true }
	oldGOOS := runtime.GOOS

	delete(clear, oldGOOS)
	Clear()
	if called {
		t.Error("Clear should not call unknown OS function")
	}
}

func TestClear_ValueCalled(t *testing.T) {
	called := false
	clear["test"] = func() { called = true }
	originalGetGOOS := getGOOS
	getGOOS = func() string { return "test" }
	defer func() { getGOOS = originalGetGOOS }()

	Clear()
	if !called {
		t.Error("Clear() should call the function for 'test' OS")
	}
}
