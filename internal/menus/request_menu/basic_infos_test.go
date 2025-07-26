package request_menu

import (
	"net/http"
	"strings"
	"testing"

	config_module "github.com/diogopereiradev/httpzen/internal/config"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/diogopereiradev/httpzen/internal/utils/http_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/stretchr/testify/assert"
)

func Test_basic_infos_Render(t *testing.T) {
	resp := &request_module.RequestResponse{
		ExecutionTime: 100,
		HttpVersion:   "HTTP/1.1",
		Method:        "GET",
		StatusMessage: "200 OK",
		Request:       request_module.RequestOptions{Url: "http://localhost"},
		Result:        "result",
		Body: []http_utility.HttpContentData{
			{ContentType: "application/json", Value: "{}", Key: ""},
		},
	}

	config := &config_module.Config{SlowResponseThreshold: 50}
	m := &Model{response: resp, config: config}

	res := basic_infos_Render(m)
	if !strings.Contains(res, "HTTP/1.1 GET 200 OK") {
		t.Errorf("expected status line in output")
	}

	if !strings.Contains(res, "Response Time:") {
		t.Errorf("expected response time in output")
	}
}

func Test_basic_infos_Render_NoBody(t *testing.T) {
	resp := &request_module.RequestResponse{
		ExecutionTime: 100,
		HttpVersion:   "HTTP/1.1",
		Method:        "GET",
		StatusMessage: "200 OK",
		Request:       request_module.RequestOptions{Url: "http://localhost"},
		Result:        "result",
	}

	config := &config_module.Config{SlowResponseThreshold: 50}
	m := &Model{response: resp, config: config}

	res := basic_infos_Render(m)
	if !strings.Contains(res, "No request body available.") {
		t.Errorf("expected no body message in output")
	}
}

func Test_basic_infos_Render_YellowSlowResponse(t *testing.T) {
	resp := &request_module.RequestResponse{
		ExecutionTime: 200,
		HttpVersion:   "HTTP/1.1",
		Method:        "GET",
		StatusMessage: "500 Internal Server Error",
		Request:       request_module.RequestOptions{Url: "http://localhost"},
		Result:        "result",
	}

	config := &config_module.Config{SlowResponseThreshold: 250}
	m := &Model{response: resp, config: config}

	res := basic_infos_Render(m)
	if !strings.Contains(res, "200.00ms (slow)") {
		t.Errorf("expected slow response time in output")
	}
}

func Test_basic_infos_Render_Paged(t *testing.T) {
	resp := &request_module.RequestResponse{
		ExecutionTime: 10,
		HttpVersion:   "HTTP/2",
		Method:        "POST",
		StatusMessage: "201 Created",
		Request:       request_module.RequestOptions{Url: "http://test"},
		Result:        strings.Repeat("a\n", 100),
	}

	config := &config_module.Config{SlowResponseThreshold: 50}
	m := &Model{response: resp, config: config}

	res := basic_infos_Render_Paged(m)
	if !strings.Contains(res, "HTTP/2 POST 201 Created") {
		t.Errorf("expected status line in paged output")
	}
}

func Test_basic_infos_Render_Paged_LinesGreaterThanMaxLines(t *testing.T) {
	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 16, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	resp := &request_module.RequestResponse{
		ExecutionTime: 10,
		HttpVersion:   "HTTP/2",
		Method:        "POST",
		StatusMessage: "201 Created",
		Request:       request_module.RequestOptions{Url: "http://test"},
		Result:        strings.Repeat("a\n", 100),
	}

	config := &config_module.Config{SlowResponseThreshold: 50}
	m := &Model{response: resp, config: config}

	
	res := basic_infos_Render_Paged(m)
	if !strings.Contains(res, "[1-0/9 lines]") {
		t.Errorf("expected paged output with line count")
	}
}

func Test_basic_infos_ScrollUp(t *testing.T) {
	resp := &request_module.RequestResponse{}
	config := &config_module.Config{}

	m := &Model{response: resp, config: config, infosScrollOffset: 2}
	basic_infos_ScrollUp(m)

	if m.infosScrollOffset != 1 {
		t.Errorf("ScrollUp did not decrement offset")
	}

	m.infosScrollOffset = 0
	basic_infos_ScrollUp(m)

	if m.infosScrollOffset != 0 {
		t.Errorf("ScrollUp should not go below zero")
	}
}

func Test_basic_infos_ScrollDown(t *testing.T) {
	resp := &request_module.RequestResponse{}
	config := &config_module.Config{}

	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	m := &Model{response: resp, config: config, infosScrollOffset: 0}
	m.infosLinesAmount = 3
	basic_infos_ScrollDown(m)
	assert.Equal(t, 1, m.infosScrollOffset)
	basic_infos_ScrollDown(m)
	assert.Equal(t, 1, m.infosScrollOffset)
}

func Test_basic_infos_ScrollDown_Borders(t *testing.T) {
	resp := &request_module.RequestResponse{}
	config := &config_module.Config{}

	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 20, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	m := &Model{response: resp, config: config, infosScrollOffset: 0}
	m.infosLinesAmount = 0
	basic_infos_ScrollDown(m)
	assert.Equal(t, 0, m.infosScrollOffset)
	m.infosLinesAmount = 2
	basic_infos_ScrollDown(m)
	assert.Equal(t, 0, m.infosScrollOffset)
}

func Test_basic_infos_ScrollPgUp(t *testing.T) {
	resp := &request_module.RequestResponse{}
	config := &config_module.Config{}

	m := &Model{response: resp, config: config, infosScrollOffset: 6}
	basic_infos_ScrollPgUp(m)

	if m.infosScrollOffset != 1 {
		t.Errorf("ScrollPgUp did not decrement by 5")
	}

	m.infosScrollOffset = 2
	basic_infos_ScrollPgUp(m)

	if m.infosScrollOffset != 0 {
		t.Errorf("ScrollPgUp should not go below zero")
	}
}

func Test_basic_infos_ScrollPgDown(t *testing.T) {
	origGetHeightFunc := terminal_utility.GetHeightFunc
	terminal_utility.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { terminal_utility.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithReqHeaders(http.Header{"A": {"1"}, "B": {"2"}, "C": {"3"}}, 0)
	m.infosLinesAmount = 3
	basic_infos_ScrollPgDown(m)
	assert.Equal(t, 5, m.infosScrollOffset)
}

func Test_basic_infos_body_Render(t *testing.T) {
	resp := &request_module.RequestResponse{
		Body: []http_utility.HttpContentData{
			{ContentType: "application/json", Value: `{"key": "value"}`},
			{ContentType: "text/plain", Value: "plain text"},
		},
	}

	m := &Model{response: resp}
	res := basic_infos_body_Render(m)

	if !strings.Contains(res, "Request Body:") {
		t.Errorf("expected request body header in output")
	}
}

func Test_basic_infos_body_Render_TextHtmlFormatting(t *testing.T) {
	resp := &request_module.RequestResponse{
		Body: []http_utility.HttpContentData{
			{ContentType: "text/html", Value: "<div>Hello</div>"},
		},
	}

	m := &Model{response: resp}
	res := basic_infos_body_Render(m)

	if !strings.Contains(res, "Request Body:") {
		t.Errorf("expected request body header in output")
	}
}

func Test_basic_infos_body_Render_MultipartValidFilePath(t *testing.T) {
	resp := &request_module.RequestResponse{
		Body: []http_utility.HttpContentData{
			{ContentType: "multipart/form-data", Value: "./basic_infos.go"},
		},
	}

	m := &Model{response: resp}
	res := basic_infos_body_Render(m)

	if !strings.Contains(res, "Request Body:") {
		t.Errorf("expected request body header in output")
	}
}

func Test_basic_infos_body_Render_MultipartValidPathButFileNotExists(t *testing.T) {
	resp := &request_module.RequestResponse{
		Body: []http_utility.HttpContentData{
			{ContentType: "multipart/form-data", Value: "./basic_idsfnfos.go"},
		},
	}

	m := &Model{response: resp}
	res := basic_infos_body_Render(m)

	if !strings.Contains(res, "Request Body:") {
		t.Errorf("expected request body header in output")
	}
}