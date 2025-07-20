package logger_module

import (
	"bytes"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	out := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = out
	buf.ReadFrom(r)
	return buf.String()
}

func TestError(t *testing.T) {
	msg := "error message"
	output := captureOutput(func() {
		Error(msg)
	})
	if !contains(output, "An error has occurred on application execution") || !contains(output, msg) {
		t.Errorf("Error() output missing expected content")
	}
}

func TestWarn(t *testing.T) {
	msg := "warn message"
	output := captureOutput(func() {
		Warn(msg)
	})
	if !contains(output, "A warning has occurred") || !contains(output, msg) {
		t.Errorf("Warn() output missing expected content")
	}
}

func TestInfo(t *testing.T) {
	msg := "info message"
	output := captureOutput(func() {
		Info(msg)
	})
	if !contains(output, "Information") || !contains(output, msg) {
		t.Errorf("Info() output missing expected content")
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
