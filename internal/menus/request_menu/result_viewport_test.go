package request_menu

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
)

func newModelWithResult(result string, scrollOffset int) *Model {
	return &Model{
		response:           &request_module.RequestResponse{Result: result},
		resultScrollOffset: scrollOffset,
	}
}

func Test_result_viewport_Render_withResult(t *testing.T) {
	m := newModelWithResult("line1\nline2\nline3", 0)
	m.resultLinesAmount = 3
	output := result_viewport_Render(m)
	assert.Contains(t, output, "line1")
	assert.Contains(t, output, "line2")
	assert.Contains(t, output, "line3")
}

func Test_result_viewport_Render_noResult(t *testing.T) {
	m := newModelWithResult("", 0)
	output := result_viewport_Render(m)
	assert.Contains(t, output, "No response available.")
}

func Test_result_viewport_ScrollUp(t *testing.T) {
	m := newModelWithResult("a", 2)
	result_viewport_ScrollUp(m)
	assert.Equal(t, 1, m.resultScrollOffset)
	result_viewport_ScrollUp(m)
	assert.Equal(t, 0, m.resultScrollOffset)
	result_viewport_ScrollUp(m)
	assert.Equal(t, 0, m.resultScrollOffset)
}

func Test_result_viewport_ScrollDown(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 24, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithResult("a\nb\nc", 0)
	m.resultLinesAmount = 3
	result_viewport_ScrollDown(m)
	assert.Equal(t, 1, m.resultScrollOffset)
	result_viewport_ScrollDown(m)
	assert.Equal(t, 2, m.resultScrollOffset)
	result_viewport_ScrollDown(m)
	assert.Equal(t, 3, m.resultScrollOffset)
	result_viewport_ScrollDown(m)
	assert.Equal(t, 2, m.resultScrollOffset)
}

func Test_result_viewport_ScrollDown_zeroLines(t *testing.T) {
	m := newModelWithResult("", 0)
	m.resultLinesAmount = 0
	result_viewport_ScrollDown(m)
	assert.Equal(t, 0, m.resultScrollOffset)
}

func Test_result_viewport_ScrollPgUp(t *testing.T) {
	m := newModelWithResult("a\nb\nc\nd\ne\nf", 6)
	result_viewport_ScrollPgUp(m)
	assert.Equal(t, 1, m.resultScrollOffset)
	result_viewport_ScrollPgUp(m)
	assert.Equal(t, 0, m.resultScrollOffset)
}

func Test_result_viewport_ScrollPgDown(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 10, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithResult(strings.Repeat("a\n", 30), 0)
	m.resultLinesAmount = 30
	result_viewport_ScrollPgDown(m)
	assert.Equal(t, 5, m.resultScrollOffset)
}

func Test_result_viewport_ScrollPgDown_zeroLines(t *testing.T) {
	m := newModelWithResult("", 0)
	m.resultLinesAmount = 0
	result_viewport_ScrollPgDown(m)
	assert.Equal(t, 0, m.resultScrollOffset)
}

func Test_result_viewport_ScrollPgDown_maxLines(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 20, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithResult(strings.Repeat("a\n", 4), 0)
	m.resultLinesAmount = 4
	result_viewport_ScrollPgDown(m)
	assert.Equal(t, 0, m.resultScrollOffset)
}

func Test_result_viewport_Render_json(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 24, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	jsonStr := `{"foo": "bar", "num": 1}`
	m := newModelWithResult(jsonStr, 0)
	output := result_viewport_Render(m)
	assert.Contains(t, output, "foo")
	assert.Contains(t, output, "bar")
}

func Test_result_viewport_Render_html(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 24, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	htmlStr := `<html><body><h1>Title</h1></body></html>`
	m := newModelWithResult(htmlStr, 0)
	output := result_viewport_Render(m)
	assert.Contains(t, output, "Title")
}

func Test_result_viewport_Render_xml(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 24, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	xmlStr := `<root><foo>bar</foo></root>`
	m := newModelWithResult(xmlStr, 0)
	output := result_viewport_Render(m)
	assert.Contains(t, output, "bar")
}

func Test_result_viewport_Render_default(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 24, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithResult("plain text", 0)
	output := result_viewport_Render(m)
	assert.Contains(t, output, "plain text")
}

func Test_result_viewport_Render_pagination(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 3 + 16, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	lines := []string{"l1", "l2", "l3", "l4", "l5"}
	m := newModelWithResult(strings.Join(lines, "\n"), 0)
	output := result_viewport_Render(m)
	assert.Contains(t, output, "[1-3/5 lines]")
}
