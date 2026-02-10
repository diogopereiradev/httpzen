package request_command

import (
	"encoding/json"
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

func isUrl(s string) bool {
	if len(s) > 1 && s[0] == ':' && s[1] >= '0' && s[1] <= '9' {
		return true
	}
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func parseInlineBody(args []string) []http_utility.HttpContentData {
	pairs := map[string]string{}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			// Remove surrounding quotes if present
			value = strings.Trim(value, "\"")
			pairs[key] = value
		}
	}

	if len(pairs) == 0 {
		return nil
	}

	jsonBytes, err := json.Marshal(pairs)
	if err != nil {
		return nil
	}

	return []http_utility.HttpContentData{
		{
			ContentType: "application/json",
			Value:       string(jsonBytes),
		},
	}
}

func Init(rootCmd *cobra.Command) {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		var method, url string
		var remainingArgs []string

		if isUrl(args[0]) {
			method = "GET"
			url = args[0]
			remainingArgs = args[1:]
		} else {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			method = http_utility.ParseHttpMethod(args[0])
			if method == "" {
				logger_module.Error(
					"Invalid HTTP method. Please provide a valid HTTP method (GET, POST, PATCH, PUT, DELETE, HEAD).",
					70,
				)
				Exit(1)
			}
			url = args[1]
			remainingArgs = args[2:]
		}

		parsedUrl := http_utility.ParseUrl(url)
		if parsedUrl == "" {
			logger_module.Error(
				"Invalid URL. Please provide a valid URL (http:// or https://).",
				70,
			)
			Exit(1)
		}

		headers, _ := cmd.Flags().GetStringSlice("header")
		flags := RequestFlags{
			Headers: headers,
			Body:    cmd.Flag("body").Value.String() == "true",
		}

		// Parse inline body from remaining args (key=value pairs)
		inlineBody := parseInlineBody(remainingArgs)
		hasInlineBody := len(inlineBody) > 0

		if (flags.Body || hasInlineBody) &&
			(method == "GET" || method == "HEAD") {
			logger_module.Error(
				"Body cannot be included in GET or HEAD requests.",
				70,
			)
			Exit(1)
			return
		}

		insecure, _ := cmd.Flags().GetBool("insecure")
		requestOptions := request_module.RequestOptions{
			Url:      parsedUrl,
			Headers:  parseHeaders(flags.Headers),
			Method:   method,
			Timeout:  30 * time.Second,
			Insecure: insecure,
		}

		var body []http_utility.HttpContentData
		if hasInlineBody {
			body = inlineBody
		} else if flags.Body {
			BodyMenuNewFunc(&requestOptions, &body)
		}
		requestOptions.Body = body

		res := RunRequestFunc(requestOptions)
		RequestMenuNewFunc(&res)
	}

	rootCmd.Flags().
		BoolP("body", "b", false, "Include body in the request (default: false)")
	rootCmd.Flags().
		StringSliceP("header", "H", []string{}, "Add a header to the request (can be used multiple times)")
	rootCmd.Flags().
		BoolP("insecure", "k", false, "Allow insecure SSL certificates (Self-Signed)")
}
