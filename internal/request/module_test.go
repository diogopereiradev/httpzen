package request_module

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	ip_cache_module "github.com/diogopereiradev/httpzen/internal/cache"
	"github.com/go-resty/resty/v2"
)

func Test_handleHttpMethod(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"Get", "Get"},
		{"Post", "Post"},
		{"Put", "Put"},
		{"Delete", "Delete"},
		{"Patch", "Patch"},
		{"Head", "Head"},
		{"Unknown", "GET"},
		{"", "GET"},
	}
	for _, c := range cases {
		if got := handleHttpMethod(c.input); got != c.expected {
			t.Errorf("handleHttpMethod(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

// URL Handling Tests
func Test_handleUrl_valid(t *testing.T) {
	url := "https://example.com"
	got := handleUrl(url)
	if got != url {
		t.Errorf("handleUrl(%q) = %q, want %q", url, got, url)
	}
}

func Test_handleUrl_invalid(t *testing.T) {
	called := false
	Exit = func(code int) { called = true }
	defer func() { Exit = os.Exit }()
	_ = handleUrl("ftp://example.com")
	if !called {
		t.Error("handleUrl should call Exit for invalid URL")
	}
}

// Execution Time Tests
func Test_handleExecutionTimeInMilliseconds(t *testing.T) {
	start := time.Now().Add(-1500 * time.Millisecond)
	ms := handleExecutionTimeInMilliseconds(start)
	if ms < 1400 || ms > 1600 {
		t.Errorf("handleExecutionTimeInMilliseconds returned %f, want ~1500", ms)
	}
}

// Domain IPs Lookup Tests
func Test_handleDomainIpsLookup_empty(t *testing.T) {
	res := &resty.Response{}
	res.Request = &resty.Request{RawRequest: &http.Request{URL: &url.URL{Host: "invalidhost"}}}
	ips := handleDomainIpsLookup(res)
	if len(ips) != 0 {
		t.Errorf("handleDomainIpsLookup should return empty slice for invalid host")
	}
}

func Test_handleDomainIpsLookup_with_port(t *testing.T) {
	res := &resty.Response{}
	res.Request = &resty.Request{RawRequest: &http.Request{URL: &url.URL{Host: "localhost:8080"}}}
	_ = handleDomainIpsLookup(res)
}

func Test_handleDomainIpsLookup_fetchIpInfo_error(t *testing.T) {
	res := &resty.Response{}
	res.Request = &resty.Request{RawRequest: &http.Request{URL: &url.URL{Host: "localhost"}}}

	original := fetchIpInfoFunc
	fetchIpInfoFunc = func(ipType, ip string) (IpInfo, error) {
		return IpInfo{}, fmt.Errorf("mock error")
	}
	defer func() { fetchIpInfoFunc = original }()

	ips := handleDomainIpsLookup(res)
	if len(ips) != 0 {
		t.Errorf("handleDomainIpsLookup should return empty slice when fetchIpInfo returns error, got: %v", ips)
	}
}

// fetchIpInfo Tests
func Test_fetchIpInfo_cache(t *testing.T) {
	key := "127_0_0_1"

	ip_cache_module.SetIpInfoToCache(key, map[string]any{"Type": "IPv4", "Ip": "127.0.0.1"})

	info, err := fetchIpInfo("IPv4", "127.0.0.1")
	if err != nil {
		t.Errorf("fetchIpInfo returned error: %v", err)
	}

	if info.Ip != "127.0.0.1" {
		t.Errorf("fetchIpInfo returned wrong Ip: %v", info.Ip)
	}
}

// HandleRequest Tests
func Test_HandleRequest_invalid_url(t *testing.T) {
	called := false
	Exit = func(code int) { called = true }
	defer func() { Exit = os.Exit }()

	opts := RequestOptions{Url: "invalid", Method: "Get"}
	_ = HandleRequest(opts)
	if !called {
		t.Error("HandleRequest should call Exit for invalid URL")
	}
}

func Test_HandleRequest_valid(t *testing.T) {
	opts := RequestOptions{
		Timeout: 2 * time.Second,
		Headers: map[string]string{"Accept": "application/json"},
		Cookies: map[string]string{"testcookie": "testvalue"},
		Url:     "https://google.com",
		Method:  "GET",
	}

	resp := HandleRequest(opts)
	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		t.Errorf("HandleRequest returned wrong status code: %v", resp.StatusCode)
	}

	if resp.Method != "GET" {
		t.Errorf("HandleRequest returned wrong method: %v", resp.Method)
	}

	if resp.Host == "" {
		t.Errorf("HandleRequest returned empty host")
	}
}
