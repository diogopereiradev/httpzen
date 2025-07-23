package request_menu

import (
	"strings"
	"testing"
)

func Test_navigation_options_Render(t *testing.T) {
	result := navigation_options_Render()

	if !strings.Contains(result, "Use left/right arrows to navigate between tabs") {
		t.Errorf("Expected navigation instructions in output, got: %s", result)
	}
	if !strings.Contains(result, "'q' to quit") {
		t.Errorf("Expected quit instruction in output, got: %s", result)
	}
	if !strings.Contains(result, "'c' to copy response") {
		t.Errorf("Expected copy response instruction in output, got: %s", result)
	}
	if !strings.Contains(result, "'b' to benchmark") {
		t.Errorf("Expected benchmark instruction in output, got: %s", result)
	}
	if !strings.Contains(result, "'r' to resend request") {
		t.Errorf("Expected resend request instruction in output, got: %s", result)
	}
}
