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

func TestGetLogWidth(t *testing.T) {
	tests := []struct {
		maxWidth int
		expected int
	}{
		{90, 90},
		{0, 90},
	}

	for _, test := range tests {
		result := getLogWidth(test.maxWidth)
		if result != test.expected {
			t.Errorf("getLogWidth(%d) = %d; want %d", test.maxWidth, result, test.expected)
		}
	}
}

func TestError(t *testing.T) {
	msg := "error message"
	output := captureOutput(func() {
		Error(msg, 90)
	})
	if !contains(output, "An error has occurred on application execution") || !contains(output, msg) {
		t.Errorf("Error() output missing expected content")
	}
}

func TestWarn(t *testing.T) {
	msg := "warn message"
	output := captureOutput(func() {
		Warn(msg, 90)
	})
	if !contains(output, "A warning has occurred") || !contains(output, msg) {
		t.Errorf("Warn() output missing expected content")
	}
}

func TestInfo(t *testing.T) {
	msg := "info message"
	output := captureOutput(func() {
		Info(msg, 90)
	})
	if !contains(output, "Information") || !contains(output, msg) {
		t.Errorf("Info() output missing expected content")
	}
}

func TestSuccess(t *testing.T) {
	msg := "success message"
	output := captureOutput(func() {
		Success(msg, 90)
	})
	if !contains(output, "Action Succeeded") || !contains(output, msg) {
		t.Errorf("Success() output missing expected content")
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
