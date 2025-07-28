package ip_utility

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	ip_cache_module "github.com/diogopereiradev/httpzen/internal/cache"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestCache() func() {
	configPath := "/tmp/httpzen_test_cache"
	os.MkdirAll(configPath, 0755)
	cacheFile := filepath.Join(configPath, "ip_test_cache.json")
	os.Remove(cacheFile)
	return func() {
		os.Remove(cacheFile)
	}
}

func TestFetchIpInfo_CacheHit(t *testing.T) {
	teardown := setupTestCache()
	defer teardown()

	info := map[string]any{
		"Type":      "IPv4",
		"Ip":        "127.0.0.1",
		"Decimal":   "12345",
		"Hostname":  "localhost",
		"ASN":       "AS123",
		"ISP":       "ISP",
		"City":      "City",
		"Country":   "Country",
		"State":     "State",
		"Latitude":  1.23,
		"Longitude": 4.56,
	}
	ip_cache_module.SetIpInfoToCache("127.0.0.1", info)

	result, err := FetchIpInfo("IPv4", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", result.Ip)
	assert.Equal(t, "IPv4", result.Type)
	assert.Equal(t, "12345", result.Decimal)
	assert.Equal(t, "localhost", result.Hostname)
	assert.Equal(t, "AS123", result.ASN)
	assert.Equal(t, "ISP", result.ISP)
	assert.Equal(t, "City", result.City)
	assert.Equal(t, "Country", result.Country)
	assert.Equal(t, "State", result.State)
	assert.Equal(t, 1.23, result.Latitude)
	assert.Equal(t, 4.56, result.Longitude)
}

func TestFetchIpInfo_Api(t *testing.T) {
	defer ip_cache_module.ClearCache()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"hostname": "host",
			"postal":   "99999",
			"org":      "ASORG",
			"city":     "C",
			"region":   "R",
			"country":  "BR",
			"loc":      "-10.1,20.2",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	oldApiUrl := IpInfoApiUrlFormat
	IpInfoApiUrlFormat = ts.URL + "/%s/json"
	defer func() { IpInfoApiUrlFormat = oldApiUrl }()

	info, err := FetchIpInfo("IPv4", "testip")
	assert.NoError(t, err)
	assert.Equal(t, "testip", info.Ip)
	assert.Equal(t, "IPv4", info.Type)
	assert.Equal(t, "99999", info.Decimal)
	assert.Equal(t, "host", info.Hostname)
	assert.Equal(t, "ASORG", info.ASN)
	assert.Equal(t, "ASORG", info.ISP)
	assert.Equal(t, "C", info.City)
	assert.Equal(t, "BR", info.Country)
	assert.Equal(t, "R", info.State)
	assert.Equal(t, -10.1, info.Latitude)
	assert.Equal(t, 20.2, info.Longitude)
}

func TestLookupDomainIps(t *testing.T) {
	res := &resty.Response{
		Request: &resty.Request{
			RawRequest: &http.Request{
				URL: &url.URL{Host: "localhost"},
			},
		},
	}
	ips := LookupDomainIps(res)
	assert.True(t, len(ips) >= 1)
	for _, ip := range ips {
		assert.NotEmpty(t, ip.Ip)
		assert.True(t, ip.Type == "IPv4" || ip.Type == "IPv6")
	}
}

func TestLookupDomainIps_Error(t *testing.T) {
	res := &resty.Response{
		Request: &resty.Request{
			RawRequest: &http.Request{
				URL: &url.URL{Host: "notfound"},
			},
		},
	}
	ips := LookupDomainIps(res)
	assert.Len(t, ips, 0)
}

func TestLookupDomainIps_HostnameWithPort(t *testing.T) {
	res := &resty.Response{
		Request: &resty.Request{
			RawRequest: &http.Request{
				URL: &url.URL{Host: "localhost:8080"},
			},
		},
	}
	ips := LookupDomainIps(res)
	assert.True(t, len(ips) >= 1)
	for _, ip := range ips {
		assert.NotEmpty(t, ip.Ip)
		assert.True(t, ip.Type == "IPv4" || ip.Type == "IPv6")
	}
}

func TestLookupDomainIps_Ipv6Protocol(t *testing.T) {
	originalLookupIP := IpLookupFunc
	defer func() { IpLookupFunc = originalLookupIP }()
	IpLookupFunc = func(host string) ([]net.IP, error) {
		return []net.IP{net.ParseIP("::1")}, nil
	}

	res := &resty.Response{
		Request: &resty.Request{
			RawRequest: &http.Request{
				URL: &url.URL{Host: "localhost"},
			},
		},
	}
	ips := LookupDomainIps(res)
	assert.Len(t, ips, 1)
	assert.Equal(t, "::1", ips[0].Ip)
	assert.Equal(t, "IPv6", ips[0].Type)
}
