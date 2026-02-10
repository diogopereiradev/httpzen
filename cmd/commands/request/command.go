package request_command

import (
	"net/http"
	"os"
	"strings"
	"time"

	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	"github.com/diogopereiradev/httpzen/internal/menus/body_menu"
	"github.com/diogopereiradev/httpzen/internal/menus/request_menu"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/diogopereiradev/httpzen/internal/utils/http_utility"
	"github.com/spf13/cobra"
)

type RequestFlags struct {
	Headers []string
	Body    bool
}

var Exit = os.Exit
var RunRequestFunc = request_module.RunRequest
var BodyMenuNewFunc = body_menu.New
var RequestMenuNewFunc = request_menu.New

func parseHeaders(headers []string) http.Header {
	result := http.Header{}
	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result.Add(key, value)
		}
	}
	return result
}

func parseMethodUrlAndRemainingArgs(cmd *cobra.Command, args []string) (method string, url string, remainingArgs []string, ok bool) {
	if len(args) < 1 {
		cmd.Help()
		return "", "", nil, false
	}

	if http_utility.CheckIsUrl(args[0]) {
		return "GET", args[0], args[1:], true
	}

	if len(args) < 2 {
		cmd.Help()
		return "", "", nil, false
	}

	method = http_utility.ParseHttpMethod(args[0])
	if method == "" {
		logger_module.Error(
			"Invalid HTTP method. Please provide a valid HTTP method (GET, POST, PATCH, PUT, DELETE, HEAD).",
			70,
		)
		Exit(1)
		return "", "", nil, false
	}

	return method, args[1], args[2:], true
}

func parseAndValidateUrl(url string) string {
	parsedUrl := http_utility.ParseUrl(url)
	if parsedUrl == "" {
		logger_module.Error(
			"Invalid URL. Please provide a valid URL (http:// or https://).",
			70,
		)
		Exit(1)
		return ""
	}
	return parsedUrl
}

func getRequestFlags(cmd *cobra.Command) RequestFlags {
	headers, _ := cmd.Flags().GetStringSlice("header")
	body, _ := cmd.Flags().GetBool("body")

	return RequestFlags{
		Headers: headers,
		Body:    body,
	}
}

func validateBodyAllowed(method string, wantsBody bool) bool {
	if !wantsBody {
		return true
	}
	if method == "GET" || method == "HEAD" {
		logger_module.Error(
			"Body cannot be included in GET or HEAD requests.",
			70,
		)
		Exit(1)
		return false
	}
	return true
}

func buildRequestOptions(cmd *cobra.Command, method string, parsedUrl string, flags RequestFlags) request_module.RequestOptions {
	insecure, _ := cmd.Flags().GetBool("insecure")
	return request_module.RequestOptions{
		Url:      parsedUrl,
		Headers:  parseHeaders(flags.Headers),
		Method:   method,
		Timeout:  30 * time.Second,
		Insecure: insecure,
	}
}

func resolveRequestBody(requestOptions *request_module.RequestOptions, flags RequestFlags, inlineBody []http_utility.HttpContentData) []http_utility.HttpContentData {
	if len(inlineBody) > 0 {
		return inlineBody
	}

	var body []http_utility.HttpContentData
	if flags.Body {
		BodyMenuNewFunc(requestOptions, &body)
	}
	return body
}

func registerFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("body", "b", false, "Include body in the request (default: false)")
	cmd.Flags().StringSliceP("header", "H", []string{}, "Add a header to the request (can be used multiple times)")
	cmd.Flags().BoolP("insecure", "k", false, "Allow insecure SSL certificates (Self-Signed)")
}

func Init(rootCmd *cobra.Command) {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		method, url, remainingArgs, ok := parseMethodUrlAndRemainingArgs(cmd, args)
		if !ok {
			return
		}

		parsedUrl := parseAndValidateUrl(url)
		flags := getRequestFlags(cmd)

		inlineBody := http_utility.ParseInlineBody(remainingArgs)
		wantsBody := flags.Body || len(inlineBody) > 0
		if !validateBodyAllowed(method, wantsBody) {
			return
		}

		requestOptions := buildRequestOptions(cmd, method, parsedUrl, flags)
		requestOptions.Body = resolveRequestBody(&requestOptions, flags, inlineBody)

		res := RunRequestFunc(requestOptions)
		RequestMenuNewFunc(&res)
	}
	registerFlags(rootCmd)
}
