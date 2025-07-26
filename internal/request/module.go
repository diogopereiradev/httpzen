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
	Timeout time.Duration                  `json:"timeout"`
	Headers http.Header                    `json:"headers"`
	Body    []http_utility.HttpContentData `json:"body"`
	Url     string                         `json:"url"`
	Method  string                         `json:"method"`
}

var Exit = os.Exit

var restyNew = resty.New
var parseHttpMethod = http_utility.ParseHttpMethod
var parseUrl = http_utility.ParseUrl
var parseExecutionTimeInMilliseconds = http_utility.ParseExecutionTimeInMilliseconds
var getConfig = config_module.GetConfig
var lookupDomainIps = ip_utility.LookupDomainIps
var loggerError = logger_module.Error
var parseApplicationJson = http_utility.ParseApplicationJson
var parseMultipartFormData = http_utility.ParseMultipartFormData
var parseUrlEncodedForm = http_utility.ParseUrlEncodedForm

type RequestResponse struct {
	HttpVersion   string                         `json:"http_version"`
	StatusMessage string                         `json:"status_message"`
	StatusCode    int                            `json:"status_code"`
	ExecutionTime float64                        `json:"execution_time"`
	Headers       http.Header                    `json:"headers"`
	Body          []http_utility.HttpContentData `json:"body"`
	Cookies       []*http.Cookie                 `json:"cookies"`
	Request       RequestOptions                 `json:"request"`
	Path          string                         `json:"path"`
	Host          string                         `json:"host"`
	Method        string                         `json:"method"`
	IpInfos       []ip_utility.LookupIpInfo      `json:"ip_infos"`
	SlowResponse  bool                           `json:"slow_response"`
	Result        string                         `json:"result"`
}

func RunRequest(options RequestOptions) RequestResponse {
	method := parseHttpMethod(options.Method)
	url := parseUrl(options.Url)
	if url == "" {
		loggerError("Invalid URL. Please provide a valid URL (http:// or https://).", 70)
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
		loggerError("Failed to execute HTTP request: " + err.Error(), 70)
		Exit(1)
		return RequestResponse{}
	}

	executionTime := parseExecutionTimeInMilliseconds(startTime)
	config := getConfig()

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
		IpInfos:       lookupDomainIps(res),
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
		return parseApplicationJson(body[0])
	}

	if contentType == "multipart/form-data" {
		return parseMultipartFormData(body)
	}

	if contentType == "application/x-www-form-urlencoded" {
		return parseUrlEncodedForm(body)
	}

	return http_utility.HandleParseResult{
		ContentTypeHeader: contentType,
		Result:            body[0].Value,
	}
}
