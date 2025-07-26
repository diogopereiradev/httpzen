package request_module

import (
	"net/http"
	"os"
	"time"

	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	"github.com/diogopereiradev/httpzen/internal/utils/http_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/ip_utility"
	"github.com/go-resty/resty/v2"
)

type RequestOptions struct {
	Timeout time.Duration                         `json:"timeout"`
	Headers http.Header                           `json:"headers"`
	Body    []http_utility.HttpContentData        `json:"body"`
	Url     string                                `json:"url"`
	Method  string                                `json:"method"`
}

type RequestResponse struct {
	HttpVersion   string                                 `json:"http_version"`
	StatusMessage string                                 `json:"status_message"`
	StatusCode    int                                    `json:"status_code"`
	ExecutionTime float64                                `json:"execution_time"`
	Headers       http.Header                            `json:"headers"`
	Body          []http_utility.HttpContentData         `json:"body"`
	Cookies       []*http.Cookie                         `json:"cookies"`
	Request       RequestOptions                         `json:"request"`
	Path          string                                 `json:"path"`
	Host          string                                 `json:"host"`
	Method        string                                 `json:"method"`
	IpInfos       []ip_utility.LookupIpInfo              `json:"ip_infos"`
	SlowResponse  bool                                   `json:"slow_response"`
	Result        string                                 `json:"result"`
}

var Exit = os.Exit
var restyNew = resty.New

func RunRequest(options RequestOptions) RequestResponse {
	method := http_utility.ParseHttpMethod(options.Method)
	url := http_utility.ParseUrl(options.Url)
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

	reqBody := HandleBody(options.Body)
	if reqBody.ContentTypeHeader != "" {
		headers["Content-Type"] = reqBody.ContentTypeHeader
		options.Headers.Set("Content-Type", reqBody.ContentTypeHeader)
	}

	req.SetHeaders(headers)
	req.SetBody(reqBody.Result)

	startTime := time.Now()

	res, err := req.Execute(method, url)
	if err != nil {
		logger_module.Error("Failed to execute HTTP request: " + err.Error())
		Exit(1)
		return RequestResponse{}
	}

	executionTime := http_utility.ParseExecutionTimeInMilliseconds(startTime)
	config := config_module.GetConfig()

	return RequestResponse{
		HttpVersion:   res.RawResponse.Proto,
		Result:        res.String(),
		StatusMessage: res.Status(),
		StatusCode:    res.StatusCode(),
		ExecutionTime: executionTime,
		Headers:       res.Header(),
		Body:          options.Body,
		Cookies:       res.Cookies(),
		Path:          res.Request.RawRequest.URL.Path,
		Host:          res.Request.RawRequest.URL.Host,
		Method:        res.Request.Method,
		IpInfos:       ip_utility.LookupDomainIps(res),
		SlowResponse:  executionTime > float64(config.SlowResponseThreshold),
		Request: RequestOptions{
			Url:     url,
			Headers: options.Headers,
			Method:  method,
			Timeout: options.Timeout,
			Body:    options.Body,
		},
	}
}

func HandleBody(body []http_utility.HttpContentData) http_utility.HandleParseResult {
	if len(body) == 0 {
		return http_utility.HandleParseResult{}
	}

	contentType := body[0].ContentType

	if contentType == "application/json" {
		return http_utility.ParseApplicationJson(body[0])
	}

	if contentType == "multipart/form-data" {
		return http_utility.ParseMultipartFormData(body)
	}

	if contentType == "application/x-www-form-urlencoded" {
		return http_utility.ParseUrlEncodedForm(body)
	}

	return http_utility.HandleParseResult{
		ContentTypeHeader: contentType,
		Result:            body[0].Value,
	}
}
