package content_type

import (
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"json object", `{"key": "value"}`, "json"},
		{"json array", `[1,2,3]`, "json"},
		{"invalid json", `{key: value}`, "text"},
		{"html doctype", `<!DOCTYPE html><html></html>`, "html"},
		{"html tag", `<html></html>`, "html"},
		{"xml declaration", `<?xml version=\"1.0\"?><root></root>`, "xml"},
		{"xml tag", `<root></root>`, "xml"},
		{"xml single tag", `<root>`, "xml"},
		{"xml with spaces", `   <root></root>   `, "xml"},
		{"plain text", `hello world`, "text"},
		{"empty string", `   `, "text"},
		{"json with spaces", `   {"a":1}   `, "json"},
		{"html with spaces", `   <html>`, "html"},
		{"xml with suffix", `<root>`, "xml"},
		{"angle brackets only", `<notxml`, "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Detect(tt.input); got != tt.expect {
				t.Errorf("Detect(%q) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}
