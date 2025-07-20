package request_module

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	ip_cache_module "github.com/diogopereiradev/httpzen/internal/cache"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	"github.com/go-resty/resty/v2"
)

type IpInfo struct {
	Type      string  `json:"type"`
	Ip        string  `json:"ip"`
	Decimal   string  `json:"decimal,omitempty"`
	Hostname  string  `json:"hostname,omitempty"`
	ASN       string  `json:"asn,omitempty"`
	ISP       string  `json:"isp,omitempty"`
	State     string  `json:"state,omitempty"`
	City      string  `json:"city,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type RequestOptions struct {
	Timeout      time.Duration     `json:"timeout"`
	Headers      map[string]string `json:"headers"`
	Body         string            `json:"body"`
	Cookies      map[string]string `json:"cookies"`
	Url          string            `json:"url"`
	Method       string            `json:"method"`
	SlowResponse bool              `json:"slow_response"`
}

type RequestResponse struct {
	StatusCode    int            `json:"status_code"`
	ExecutionTime string         `json:"execution_time"`
	Headers       http.Header    `json:"headers"`
	Body          string         `json:"body"`
	Cookies       []*http.Cookie `json:"cookies"`
	Path          string         `json:"path"`
	Host          string         `json:"host"`
	Method        string         `json:"method"`
	IpInfos       []IpInfo       `json:"ip_infos"`
	SlowResponse  bool           `json:"slow_response"`
	Result        string         `json:"result"`
}

var Exit = os.Exit

var fetchIpInfoFunc = fetchIpInfo

func HandleRequest(options RequestOptions) RequestResponse {
	method := handleHttpMethod(options.Method)
	url := handleUrl(options.Url)

	client := resty.New()
	client.SetTimeout(options.Timeout)

	req := client.R()
	req.SetHeaders(options.Headers)
	req.SetBody(options.Body)

	var cookies []*http.Cookie
	for k, v := range options.Cookies {
		cookies = append(cookies, &http.Cookie{Name: k, Value: v})
	}
	if len(cookies) > 0 {
		req.SetCookies(cookies)
	}

	startTime := time.Now()

	res, err := req.Execute(method, url)
	if err != nil {
		logger_module.Error("Failed to execute HTTP request: " + err.Error())
		Exit(1)
	}

	executionTime := handleExecutionTimeInMilliseconds(startTime)
	config := config_module.GetConfig()

	return RequestResponse{
		Result:        res.String(),
		StatusCode:    res.StatusCode(),
		ExecutionTime: fmt.Sprintf("%.2f", executionTime),
		Headers:       res.Header(),
		Body:          options.Body,
		Cookies:       res.Cookies(),
		Path:          res.Request.RawRequest.URL.Path,
		Host:          res.Request.RawRequest.URL.Host,
		Method:        res.Request.Method,
		IpInfos:       handleDomainIpsLookup(res),
		SlowResponse:  executionTime > float64(config.SlowResponseThreshold),
	}
}

func handleHttpMethod(method string) string {
	switch method {
	case "Get", "Post", "Put", "Delete", "Patch", "Head":
		return method
	default:
		return "GET"
	}
}

func handleUrl(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		logger_module.Error("Please enter a valid URL: " + url)
		Exit(1)
	}
	return url
}

func handleExecutionTimeInMilliseconds(start time.Time) float64 {
	executionTime := time.Since(start)
	ms := float64(executionTime.Nanoseconds()) / 1e6
	return ms
}

func handleDomainIpsLookup(res *resty.Response) []IpInfo {
	host := res.Request.RawRequest.URL.Host
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return []IpInfo{}
	}

	var ipList []IpInfo
	for _, ip := range ips {
		ipType := ""
		if ip.To4() != nil {
			ipType = "IPv4"
		} else {
			ipType = "IPv6"
		}

		ipInfo, err := fetchIpInfoFunc(ipType, ip.String())
		if err != nil {
			continue
		}

		ipList = append(ipList, ipInfo)
	}
	return ipList
}

func fetchIpInfo(ipType string, ip string) (IpInfo, error) {
	key := strings.ReplaceAll(ip, ".", "_")
	if cached, ok := ip_cache_module.GetIpInfoFromCache(key); ok {
		var info IpInfo
		b, _ := json.Marshal(cached)
		_ = json.Unmarshal(b, &info)
		return info, nil
	}

	apiUrl := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	client := resty.New()
	res, _ := client.R().Get(apiUrl)

	type ipinfoResponse struct {
		Hostname string `json:"hostname"`
		Postal   string `json:"postal"`
		Org      string `json:"org"`
		City     string `json:"city"`
		Region   string `json:"region"`
		Country  string `json:"country"`
		Loc      string `json:"loc"`
		ASN      string `json:"asn"`
		ISP      string `json:"isp"`
	}

	var resp ipinfoResponse
	json.Unmarshal(res.Body(), &resp)

	latitude := 0.0
	longitude := 0.0

	if resp.Loc != "" {
		parts := strings.Split(resp.Loc, ",")
		if len(parts) == 2 {
			fmt.Sscanf(parts[0], "%f", &latitude)
			fmt.Sscanf(parts[1], "%f", &longitude)
		}
	}

	info := IpInfo{
		Type:      ipType,
		Ip:        ip,
		Decimal:   resp.Postal,
		Hostname:  resp.Hostname,
		ASN:       resp.Org,
		ISP:       resp.Org,
		City:      resp.City,
		State:     resp.Region,
		Latitude:  latitude,
		Longitude: longitude,
	}

	ip_cache_module.SetIpInfoToCache(key, map[string]any{
		"Type":      info.Type,
		"Ip":        info.Ip,
		"Decimal":   info.Decimal,
		"Hostname":  info.Hostname,
		"ASN":       info.ASN,
		"ISP":       info.ISP,
		"City":      info.City,
		"State":     info.State,
		"Latitude":  info.Latitude,
		"Longitude": info.Longitude,
	})
	return info, nil
}
