package request_module

import (
	"errors"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	config_module "github.com/diogopereiradev/httpzen/internal/config"
	"github.com/diogopereiradev/httpzen/internal/utils/http_utility"
	ip_utility "github.com/diogopereiradev/httpzen/internal/utils/ip_utility"
	"github.com/go-resty/resty/v2"
)

func TestRunRequest_InvalidURL(t *testing.T) {
	called := false
	Exit = func(code int) { called = true }
	loggerError = func(msg string, maxWidth int) {}
	defer func() { Exit = os.Exit; loggerError = nil }()

	options := RequestOptions{
		Url:    "invalid-url",
		Method: "GET",
	}

	resp := RunRequest(options)
	if !called {
		t.Errorf("Exit should be called on invalid URL")
	}

	if !reflect.DeepEqual(resp, RequestResponse{}) {
		t.Errorf("Expected empty response on invalid URL")
	}
}

func TestRunRequest_RequestError(t *testing.T) {
	Exit = func(code int) {}
	loggerError = func(msg string, maxWidth int) {}
	restyNew = func() *resty.Client {
		c := resty.New()
		c.OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error {
			return errors.New("fail")
		})
		return c
	}
	defer func() { Exit = os.Exit; loggerError = nil; restyNew = resty.New }()

	options := RequestOptions{
		Url:    "http://localhost",
		Method: "GET",
	}

	resp := RunRequest(options)
	if !reflect.DeepEqual(resp, RequestResponse{}) {
		t.Errorf("Expected empty response on request error")
	}
}

func TestRunRequest_Success(t *testing.T) {
	Exit = func(code int) {}
	loggerError = func(msg string, maxWidth int) {}
	getConfig = func() config_module.Config {
		return config_module.Config{SlowResponseThreshold: 1000}
	}

	lookupDomainIps = func(_ *resty.Response) []ip_utility.LookupIpInfo {
		return []ip_utility.LookupIpInfo{{Ip: "127.0.0.1"}}
	}

	defer func() {
		Exit = os.Exit
		loggerError = nil
		getConfig = config_module.GetConfig
		lookupDomainIps = ip_utility.LookupDomainIps
		restyNew = resty.New
	}()

	options := RequestOptions{
		Url:     "https://google.com",
		Method:  "GET",
		Timeout: 1 * time.Second,
		Headers: http.Header{"X-Test": {"1"}},
		Body:    []http_utility.HttpContentData{{ContentType: "text/plain", Value: "test body"}},
	}

	resp := RunRequest(options)
	if resp.Request.Url != "https://google.com" {
		t.Errorf("Expected url to be set")
	}

	if resp.Method != "GET" {
		t.Errorf("Expected method to be GET")
	}

	if len(resp.IpInfos) == 0 {
		t.Errorf("Expected IpInfos to be filled")
	} else if resp.IpInfos[0].Ip != "127.0.0.1" {
		t.Errorf("Expected fake IP info")
	}
}

func TestHandleBody_Empty(t *testing.T) {
	res := HandleBody([]http_utility.HttpContentData{})
	if !reflect.DeepEqual(res, http_utility.HandleParseResult{}) {
		t.Errorf("Expected empty result for empty body")
	}
}

func TestHandleBody_Json(t *testing.T) {
	called := false
	parseApplicationJson = func(data http_utility.HttpContentData) http_utility.HandleParseResult {
		called = true
		return http_utility.HandleParseResult{ContentTypeHeader: "application/json", Result: "{}"}
	}
	defer func() { parseApplicationJson = http_utility.ParseApplicationJson }()

	res := HandleBody([]http_utility.HttpContentData{{ContentType: "application/json", Value: "{}"}})
	if !called {
		t.Errorf("ParseApplicationJson should be called")
	}

	if res.ContentTypeHeader != "application/json" {
		t.Errorf("Expected application/json header")
	}
}

func TestHandleBody_Multipart(t *testing.T) {
	called := false
	parseMultipartFormData = func(data []http_utility.HttpContentData) http_utility.HandleParseResult {
		called = true
		return http_utility.HandleParseResult{ContentTypeHeader: "multipart/form-data", Result: "data"}
	}
	defer func() { parseMultipartFormData = http_utility.ParseMultipartFormData }()

	res := HandleBody([]http_utility.HttpContentData{{ContentType: "multipart/form-data"}})
	if !called {
		t.Errorf("ParseMultipartFormData should be called")
	}

	if res.ContentTypeHeader != "multipart/form-data" {
		t.Errorf("Expected multipart/form-data header")
	}
}

func TestHandleBody_UrlEncoded(t *testing.T) {
	called := false
	parseUrlEncodedForm = func(data []http_utility.HttpContentData) http_utility.HandleParseResult {
		called = true
		return http_utility.HandleParseResult{ContentTypeHeader: "application/x-www-form-urlencoded", Result: "foo=bar"}
	}
	defer func() { parseUrlEncodedForm = http_utility.ParseUrlEncodedForm }()

	res := HandleBody([]http_utility.HttpContentData{{ContentType: "application/x-www-form-urlencoded"}})
	if !called {
		t.Errorf("ParseUrlEncodedForm should be called")
	}

	if res.ContentTypeHeader != "application/x-www-form-urlencoded" {
		t.Errorf("Expected application/x-www-form-urlencoded header")
	}
}

func TestHandleBody_Other(t *testing.T) {
	res := HandleBody([]http_utility.HttpContentData{{ContentType: "text/plain", Value: "abc"}})
	if res.ContentTypeHeader != "text/plain" || res.Result != "abc" {
		t.Errorf("Expected passthrough for unknown content type")
	}
}
