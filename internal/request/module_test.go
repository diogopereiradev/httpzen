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

func Test_HandleHttpMethod(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"get", "GET"},
		{"post", "POST"},
		{"put", "PUT"},
		{"delete", "DELETE"},
		{"patch", "PATCH"},
		{"head", "HEAD"},
		{"", ""},
	}
	for _, c := range cases {
		if got := HandleHttpMethod(c.input); got != c.expected {
			t.Errorf("handleHttpMethod(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

// URL Handling Tests
func Test_HandleUrl_valid(t *testing.T) {
	url := "https://example.com"
	got := HandleUrl(url)
	if got != url {
		t.Errorf("handleUrl(%q) = %q, want %q", url, got, url)
	}
}

func Test_HandleUrl_invalid(t *testing.T) {
	result := HandleUrl("ftp://example.com")
	if result != "" {
		t.Errorf("handleUrl should return empty string for invalid URL, got: %q", result)
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
func Test_RunRequest_invalid_url(t *testing.T) {
	called := false
	Exit = func(code int) { called = true }
	defer func() { Exit = os.Exit }()

	opts := RequestOptions{Url: "invalid", Method: "GET"}
	_ = RunRequest(opts)
	if !called {
		t.Error("RunRequest should call Exit for invalid URL")
	}
}

func Test_RunRequest_execute_error(t *testing.T) {
	called := false
	Exit = func(code int) { called = true }
	defer func() { Exit = os.Exit }()

	opts := RequestOptions{
		Url:     "https://error.com",
		Method:  "GET",
		Timeout: 1 * time.Second,
	}

	originalRestyNew := restyNew
	restyNew = func() *resty.Client {
		client := originalRestyNew()
		client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			return fmt.Errorf("forced error")
		})
		return client
	}
	defer func() { restyNew = originalRestyNew }()

	_ = RunRequest(opts)
	if !called {
		t.Error("RunRequest should call Exit for HTTP execution error")
	}
}

func Test_RunRequest_valid(t *testing.T) {
	headers := http.Header{}
	headers.Add("Accept", "application/json")
	opts := RequestOptions{
		Timeout: 2 * time.Second,
		Headers: headers,
		Url:     "https://google.com",
		Method:  "GET",
	}

	resp := RunRequest(opts)
	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		t.Errorf("RunRequest returned wrong status code: %v", resp.StatusCode)
	}

	if resp.Method != "GET" {
		t.Errorf("RunRequest returned wrong method: %v", resp.Method)
	}

	if resp.Host == "" {
		t.Errorf("RunRequest returned empty host")
	}
}
