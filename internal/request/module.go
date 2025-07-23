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
	Country   string  `json:"country,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type RequestOptions struct {
	Timeout time.Duration `json:"timeout"`
	Headers http.Header   `json:"headers"`
	Body    string        `json:"body"`
	Url     string        `json:"url"`
	Method  string        `json:"method"`
}

type RequestResponse struct {
	HttpVersion     string         `json:"http_version"`
	StatusMessage   string         `json:"status_message"`
	StatusCode      int            `json:"status_code"`
	ExecutionTime   float64        `json:"execution_time"`
	Headers         http.Header    `json:"headers"`
	Body            string         `json:"body"`
	Cookies         []*http.Cookie `json:"cookies"`
	Request         RequestOptions `json:"request"`
	Path            string         `json:"path"`
	Host            string         `json:"host"`
	Method          string         `json:"method"`
	IpInfos         []IpInfo       `json:"ip_infos"`
	SlowResponse    bool           `json:"slow_response"`
	Result          string         `json:"result"`
}

var Exit = os.Exit

var restyNew = resty.New

var fetchIpInfoFunc = fetchIpInfo

func RunRequest(options RequestOptions) RequestResponse {
	method := HandleHttpMethod(options.Method)
	url := HandleUrl(options.Url)
	if url == "" {
		logger_module.Error("Invalid URL. Please provide a valid URL (http:// or https://).")
		Exit(1)
		return RequestResponse{}
	}

	client := restyNew()
	client.SetTimeout(options.Timeout)

	req := client.R()
	headers := make(map[string]string)
	for k, v := range options.Headers {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	req.SetHeaders(headers)
	req.SetBody(options.Body)

	startTime := time.Now()

	res, err := req.Execute(method, url)
	if err != nil {
		logger_module.Error("Failed to execute HTTP request: " + err.Error())
		Exit(1)
		return RequestResponse{}
	}

	executionTime := handleExecutionTimeInMilliseconds(startTime)
	config := config_module.GetConfig()

	return RequestResponse{
		HttpVersion:     res.RawResponse.Proto,
		Result:          res.String(),
		StatusMessage:   res.Status(),
		StatusCode:      res.StatusCode(),
		ExecutionTime:   executionTime,
		Headers:         res.Header(),
		Body:            options.Body,
		Cookies:         res.Cookies(),
		Path:            res.Request.RawRequest.URL.Path,
		Host:            res.Request.RawRequest.URL.Host,
		Method:          res.Request.Method,
		IpInfos:         handleDomainIpsLookup(res),
		SlowResponse:    executionTime > float64(config.SlowResponseThreshold),
		Request:         RequestOptions{
			Url:     url,
			Headers: options.Headers,
			Method:  method,
			Timeout: options.Timeout,
			Body:    options.Body,
		},
	}
}

func HandleHttpMethod(method string) string {
	switch strings.ToLower(method) {
	case "get":
		return "GET"
	case "post":
		return "POST"
	case "put":
		return "PUT"
	case "delete":
		return "DELETE"
	case "patch":
		return "PATCH"
	case "head":
		return "HEAD"
	default:
		return ""
	}
}

func HandleUrl(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return ""
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
		Country:   resp.Country,
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
		"Country":   info.Country,
		"State":     info.State,
		"Latitude":  info.Latitude,
		"Longitude": info.Longitude,
	})
	return info, nil
}
