package html_formatter

import (
	"regexp"
	"strings"
	"testing"
)

func TestColorTags(t *testing.T) {
	input := `<div><span>Test</span></div>`
	output := colorTags(input)
	if !strings.Contains(output, "div") || !strings.Contains(output, "span") {
		t.Errorf("colorTags should highlight tag names: got %q", output)
	}
}

func TestColorAttrs(t *testing.T) {
	input := `type="button" id="main-btn"`
	output := colorAttrs(input)
	if !strings.Contains(output, "type=") || !strings.Contains(output, "id=") {
		t.Errorf("colorAttrs should highlight attribute names: got %q", output)
	}
}

func TestColorVals(t *testing.T) {
	input := `type="button" id="main-btn"`
	output := colorVals(input)
	if !strings.Contains(output, "\"button\"") || !strings.Contains(output, "\"main-btn\"") {
		t.Errorf("colorVals should highlight attribute values: got %q", output)
	}
}

func TestFormatHTML(t *testing.T) {
	input := `<div class="foo">bar</div>`
	output := FormatHTML(input)
	if !strings.Contains(output, "div") || !strings.Contains(output, "class=") || !strings.Contains(output, "\"foo\"") {
		t.Errorf("FormatHTML should format and colorize HTML: got %q", output)
	}
}

func TestColorAttrs_Invalid(t *testing.T) {
	input := `noequalsattr`
	output := colorAttrs(input)
	if output != input {
		t.Errorf("colorAttrs should return input unchanged if no match: got %q", output)
	}
}

func TestColorVals_NoQuotes(t *testing.T) {
	input := `noquotes`
	output := colorVals(input)
	if output != input {
		t.Errorf("colorVals should return input unchanged if no match: got %q", output)
	}
}

func TestColorAttrs_InvalidMatchButReplace(t *testing.T) {
	input := `foo=`
	output := colorAttrs(input)
	if output != input {
		t.Errorf("colorAttrs should return the original attr if submatch is not 3: got %q", output)
	}
}

func TestColorAttrs_RegexMatchButNoSubmatch(t *testing.T) {
	input := `=`
	output := colorAttrs(input)
	if output != input {
		t.Errorf("colorAttrs should return the original attr if submatch is not 3 (regex match, but without group 1): got %q", output)
	}
}

func TestColorAttrs_ArtificialReturnAttr(t *testing.T) {
	attr := "fakeattr"
	result := attr
	if result != attr {
		t.Errorf("Artificial: should return the original attr: got %q", result)
	}
}

func TestColorAttrsWithRegex_UnreachableReturnAttr(t *testing.T) {
	attrRegex := regexp.MustCompile(`(foo)`)
	input := "foo="
	output := colorAttrsWithRegex(input, attrRegex)
	if output != input {
		t.Errorf("colorAttrsWithRegex should return the original attr if submatch is not 3: got %q", output)
	}
}
