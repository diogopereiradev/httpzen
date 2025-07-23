package request_command

import (
	"net/http"
	"os"
	"strings"
	"time"

	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	"github.com/diogopereiradev/httpzen/internal/menus/request_menu"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/spf13/cobra"
)

type RequestFlags struct {
	// Data
	Headers []string

	// Content types
	Json      bool
	Raw       bool
	Form      bool
	Multipart bool

	// Output
	HeadersOnly bool
	BodyOnly    bool
	MetaOnly    bool
}

var Exit = os.Exit

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

var RunRequestFunc = request_module.RunRequest
var RequestMenuNewFunc = request_menu.New

func Init(rootCmd *cobra.Command) {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}
		method := request_module.HandleHttpMethod(args[0])
		if method == "" {
			logger_module.Error("Invalid HTTP method. Please provide a valid HTTP method (GET, POST, PATCH, PUT, DELETE, HEAD).")
			Exit(1)
		}

		url := request_module.HandleUrl(args[1])
		if url == "" {
			logger_module.Error("Invalid URL. Please provide a valid URL (http:// or https://).")
			Exit(1)
		}

		headers, _ := cmd.Flags().GetStringSlice("header")
		_ = RequestFlags{
			Headers:     headers,
			Json:        cmd.Flag("json").Value.String() == "true",
			Raw:         cmd.Flag("raw").Value.String() == "true",
			Form:        cmd.Flag("form").Value.String() == "true",
			Multipart:   cmd.Flag("multipart").Value.String() == "true",
			HeadersOnly: cmd.Flag("headers").Value.String() == "true",
			BodyOnly:    cmd.Flag("body").Value.String() == "true",
			MetaOnly:    cmd.Flag("meta").Value.String() == "true",
		}

		requestOptions := request_module.RequestOptions{
			Url:     url,
			Headers: parseHeaders(headers),
			Method:  method,
			Timeout: 30 * time.Second,
		}

		res := RunRequestFunc(requestOptions)
		RequestMenuNewFunc(&res)
	}

	rootCmd.Flags().StringSliceP("header", "H", []string{}, "Add a header to the request (can be used multiple times)")

	rootCmd.Flags().BoolP("json", "j", false, "Serialize request response as JSON")
	rootCmd.Flags().BoolP("raw", "r", false, "Send raw request body")
	rootCmd.Flags().BoolP("form", "f", false, "Send form data")
	rootCmd.Flags().Bool("multipart", false, "Send multipart form data")

	rootCmd.Flags().Bool("headers", false, "Show request headers only")
	rootCmd.Flags().BoolP("body", "b", false, "Show request body only")
	rootCmd.Flags().BoolP("meta", "m", false, "Show request metadata only")
}
