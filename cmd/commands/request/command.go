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

func Init(rootCmd *cobra.Command) {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}
		method := http_utility.ParseHttpMethod(args[0])
		if method == "" {
			logger_module.Error("Invalid HTTP method. Please provide a valid HTTP method (GET, POST, PATCH, PUT, DELETE, HEAD).", 70)
			Exit(1)
		}

		url := http_utility.ParseUrl(args[1])
		if url == "" {
			logger_module.Error("Invalid URL. Please provide a valid URL (http:// or https://).", 70)
			Exit(1)
		}

		headers, _ := cmd.Flags().GetStringSlice("header")
		flags := RequestFlags{
			Headers: headers,
			Body:    cmd.Flag("body").Value.String() == "true",
		}

		if flags.Body && (method == "GET" || method == "HEAD") {
			logger_module.Error("Body cannot be included in GET or HEAD requests.", 70)
			Exit(1)
			return
		}

		requestOptions := request_module.RequestOptions{
			Url:     url,
			Headers: parseHeaders(flags.Headers),
			Method:  method,
			Timeout: 30 * time.Second,
		}

		var body []http_utility.HttpContentData
		if flags.Body {
			BodyMenuNewFunc(&requestOptions, &body)
		}
		requestOptions.Body = body

		res := RunRequestFunc(requestOptions)
		RequestMenuNewFunc(&res)
	}

	rootCmd.Flags().BoolP("body", "b", false, "Include body in the request (default: false)")
	rootCmd.Flags().StringSliceP("header", "H", []string{}, "Add a header to the request (can be used multiple times)")
}
