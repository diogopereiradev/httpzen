package request_menu

import (
	"net/http"
	"testing"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/stretchr/testify/assert"
)

func newModelWithHeaders(headers http.Header, scrollOffset int) *Model {
	return &Model{
		response:                &request_module.RequestResponse{Headers: headers},
		respHeadersScrollOffset: scrollOffset,
	}
}

func Test_response_headers_Render(t *testing.T) {
	h := http.Header{"X-Test": {"abc"}, "A": {"1"}, "B": {}}
	m := newModelWithHeaders(h, 0)
	out := response_headers_Render(m)
	assert.Contains(t, out, "X-Test")
	assert.Contains(t, out, "abc")
	assert.Contains(t, out, "A")
	assert.Contains(t, out, "1")
	assert.NotContains(t, out, "B:")
}

func Test_response_headers_Render_Paged(t *testing.T) {
	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	h := http.Header{"A": {"1"}, "B": {"2"}, "C": {"3"}}
	m := newModelWithHeaders(h, 0)
	out := response_headers_Render_Paged(m)
	assert.Contains(t, out, "A")
	assert.Contains(t, out, "B")
	assert.NotContains(t, out, "C")
	assert.Contains(t, out, "[1-2/3 lines]")

	m.respHeadersScrollOffset = 2
	out2 := response_headers_Render_Paged(m)
	assert.Contains(t, out2, "C")
	assert.Contains(t, out2, "[3-3/3 lines]")
}

func Test_response_headers_ScrollUp(t *testing.T) {
	m := newModelWithHeaders(http.Header{"A": {"1"}}, 2)
	response_headers_ScrollUp(m)
	assert.Equal(t, 1, m.respHeadersScrollOffset)
	response_headers_ScrollUp(m)
	assert.Equal(t, 0, m.respHeadersScrollOffset)
	response_headers_ScrollUp(m)
	assert.Equal(t, 0, m.respHeadersScrollOffset)
}

func Test_response_headers_ScrollDown(t *testing.T) {
	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithHeaders(http.Header{"A": {"1"}, "B": {"2"}, "C": {"3"}}, 0)
	m.respHeadersLinesAmount = 3
	response_headers_ScrollDown(m)
	assert.Equal(t, 1, m.respHeadersScrollOffset)
	response_headers_ScrollDown(m)
	assert.Equal(t, 1, m.respHeadersScrollOffset)
}

func Test_response_headers_ScrollDown_borders(t *testing.T) {
	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 20, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithHeaders(http.Header{"A": {"1"}}, 0)
	m.respHeadersLinesAmount = 0
	response_headers_ScrollDown(m)
	assert.Equal(t, 0, m.respHeadersScrollOffset)
	m.respHeadersLinesAmount = 2
	response_headers_ScrollDown(m)
	assert.Equal(t, 0, m.respHeadersScrollOffset)
}

func Test_response_headers_ScrollPgUp(t *testing.T) {
	m := newModelWithHeaders(http.Header{"A": {"1"}}, 6)
	response_headers_ScrollPgUp(m)
	assert.Equal(t, 1, m.respHeadersScrollOffset)
	response_headers_ScrollPgUp(m)
	assert.Equal(t, 0, m.respHeadersScrollOffset)
}

func Test_response_headers_ScrollPgDown(t *testing.T) {
	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithHeaders(http.Header{"A": {"1"}, "B": {"2"}, "C": {"3"}}, 0)
	m.respHeadersLinesAmount = 3
	response_headers_ScrollPgDown(m)
	assert.Equal(t, 5, m.respHeadersScrollOffset)
}

func Test_response_headers_ScrollPgDown_borders(t *testing.T) {
	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 20, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithHeaders(http.Header{"A": {"1"}}, 0)
	m.respHeadersLinesAmount = 0
	response_headers_ScrollPgDown(m)
	assert.Equal(t, 0, m.respHeadersScrollOffset)
	m.respHeadersLinesAmount = 2
	response_headers_ScrollPgDown(m)
	assert.Equal(t, 0, m.respHeadersScrollOffset)
}
