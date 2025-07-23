package json_formatter

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestFormatJSON_KeyAndString(t *testing.T) {
	input := `{"key": "value"}`
	output := FormatJSON(input)
	if !containsStyled(output, "key", lipgloss.Color("3")) {
		t.Errorf("Key not styled correctly: %v", output)
	}
	if !containsStyled(output, "value", lipgloss.Color("2")) {
		t.Errorf("String value not styled correctly: %v", output)
	}
}

func TestFormatJSON_Number(t *testing.T) {
	input := `{"num": 123.45}`
	output := FormatJSON(input)
	if !containsStyled(output, "123.45", lipgloss.Color("6")) {
		t.Errorf("Number not styled correctly: %v", output)
	}
}

func TestFormatJSON_Bool(t *testing.T) {
	input := `{"ok": true, "fail": false}`
	output := FormatJSON(input)
	if !containsStyled(output, "true", lipgloss.Color("5")) {
		t.Errorf("True not styled correctly: %v", output)
	}
	if !containsStyled(output, "false", lipgloss.Color("5")) {
		t.Errorf("False not styled correctly: %v", output)
	}
}

func TestFormatJSON_Null(t *testing.T) {
	input := `{"empty": null}`
	output := FormatJSON(input)
	if !containsStyled(output, "null", lipgloss.Color("1")) {
		t.Errorf("Null not styled correctly: %v", output)
	}
}

func TestFormatJSON_InvalidJSON(t *testing.T) {
	input := `not json`
	output := FormatJSON(input)
	if output != input {
		t.Errorf("Should return input for invalid JSON")
	}
}

func TestFormatJSON_OtherChars(t *testing.T) {
	input := `{"a":1,"b":2}`
	output := FormatJSON(input)
	if !strings.Contains(output, ",") || !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Errorf("Other chars not preserved: %v", output)
	}
}

func TestFormatJSON_EscapedString(t *testing.T) {
	input := `{"text": "line\"with\"quotes"}`
	output := FormatJSON(input)
	if !containsStyled(output, `"line\"with\"quotes"`, lipgloss.Color("2")) {
		t.Errorf("Escaped string not styled correctly: %v", output)
	}
}

func containsStyled(output, substr string, color lipgloss.Color) bool {
	style := lipgloss.NewStyle().Foreground(color)
	return strings.Contains(output, style.Render(substr))
}
