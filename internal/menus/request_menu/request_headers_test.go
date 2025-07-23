package request_menu

import (
	"net/http"
	"testing"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
	"github.com/stretchr/testify/assert"
)

func newModelWithReqHeaders(headers http.Header, scrollOffset int) *Model {
	return &Model{
		response:               &request_module.RequestResponse{Request: request_module.RequestOptions{Headers: headers}},
		reqHeadersScrollOffset: scrollOffset,
	}
}

func Test_request_headers_Render(t *testing.T) {
	h := http.Header{"X-Test": {"abc"}, "A": {"1"}, "B": {}}
	m := newModelWithReqHeaders(h, 0)
	out := request_headers_Render(m)
	assert.Contains(t, out, "X-Test")
	assert.Contains(t, out, "abc")
	assert.Contains(t, out, "A")
	assert.Contains(t, out, "1")
	assert.NotContains(t, out, "B:")
}

func Test_request_headers_Render_empty(t *testing.T) {
	m := newModelWithReqHeaders(http.Header{}, 0)
	out := request_headers_Render(m)
	assert.Contains(t, out, "No request headers found.")
}

func Test_request_headers_Render_Paged(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	h := http.Header{"A": {"1"}, "B": {"2"}, "C": {"3"}}
	m := newModelWithReqHeaders(h, 0)
	out := request_headers_Render_Paged(m)
	assert.Contains(t, out, "A")
	assert.Contains(t, out, "B")
	assert.NotContains(t, out, "C")
	assert.Contains(t, out, "[1-2/3 lines]")

	// Segunda página
	m.reqHeadersScrollOffset = 2
	out2 := request_headers_Render_Paged(m)
	assert.Contains(t, out2, "C")
	assert.Contains(t, out2, "[3-3/3 lines]")
}

func Test_request_headers_ScrollUp(t *testing.T) {
	m := newModelWithReqHeaders(http.Header{"A": {"1"}}, 2)
	request_headers_ScrollUp(m)
	assert.Equal(t, 1, m.reqHeadersScrollOffset)
	request_headers_ScrollUp(m)
	assert.Equal(t, 0, m.reqHeadersScrollOffset)
	request_headers_ScrollUp(m)
	assert.Equal(t, 0, m.reqHeadersScrollOffset)
}

func Test_request_headers_ScrollDown(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithReqHeaders(http.Header{"A": {"1"}, "B": {"2"}, "C": {"3"}}, 0)
	m.reqHeadersLinesAmount = 3
	request_headers_ScrollDown(m)
	assert.Equal(t, 1, m.reqHeadersScrollOffset)
	request_headers_ScrollDown(m)
	assert.Equal(t, 1, m.reqHeadersScrollOffset)
}

func Test_request_headers_ScrollDown_borders(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 20, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithReqHeaders(http.Header{"A": {"1"}}, 0)
	m.reqHeadersLinesAmount = 0
	request_headers_ScrollDown(m)
	assert.Equal(t, 0, m.reqHeadersScrollOffset)
	m.reqHeadersLinesAmount = 2
	request_headers_ScrollDown(m)
	assert.Equal(t, 0, m.reqHeadersScrollOffset)
}

func Test_request_headers_ScrollPgUp(t *testing.T) {
	m := newModelWithReqHeaders(http.Header{"A": {"1"}}, 6)
	request_headers_ScrollPgUp(m)
	assert.Equal(t, 1, m.reqHeadersScrollOffset)
	request_headers_ScrollPgUp(m)
	assert.Equal(t, 0, m.reqHeadersScrollOffset)
}

func Test_request_headers_ScrollPgDown(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithReqHeaders(http.Header{"A": {"1"}, "B": {"2"}, "C": {"3"}}, 0)
	m.reqHeadersLinesAmount = 3
	request_headers_ScrollPgDown(m)
	assert.Equal(t, 5, m.reqHeadersScrollOffset)
}

func Test_request_headers_ScrollPgDown_borders(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 20, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithReqHeaders(http.Header{"A": {"1"}}, 0)
	m.reqHeadersLinesAmount = 0
	request_headers_ScrollPgDown(m)
	assert.Equal(t, 0, m.reqHeadersScrollOffset)
	m.reqHeadersLinesAmount = 2
	request_headers_ScrollPgDown(m)
	assert.Equal(t, 0, m.reqHeadersScrollOffset)
}
