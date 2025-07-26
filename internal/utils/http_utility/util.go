package http_utility

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
)

type HttpContentData struct {
	ContentType string `json:"content_type"`
	Key         string `json:"key,omitempty"` // Optional key for form data/multipart
	Value       string `json:"value"`
}

type HandleParseResult struct {
	ContentTypeHeader string `json:"content_type_header"`
	Result            any    `json:"result"`
}

type FileInfoData struct {
	Name        string
	PathIsValid bool
}

func ParseApplicationJson(data HttpContentData) HandleParseResult {
	var jsonData any

	if err := json.Unmarshal([]byte(data.Value), &jsonData); err != nil {
		logger_module.Error("Failed to parse JSON body: " + err.Error())
		return HandleParseResult{}
	}
	return HandleParseResult{ContentTypeHeader: data.ContentType, Result: jsonData}
}

func ParseMultipartFormData(data []HttpContentData) HandleParseResult {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, part := range data {
		var unmarshalResult map[string]string
		var pair []string
		if err := json.Unmarshal([]byte(part.Value), &unmarshalResult); err != nil {
			pair = []string{part.Key, part.Value}
		}

		if len(pair) > 0 {
			if _, err := GetFileByPath(part.Value); err == nil {
				file, err := os.Open(part.Value)
				if err == nil {
					defer file.Close()
					if fw, err2 := writer.CreateFormFile(part.Key, part.Value); err2 == nil {
						_, _ = io.Copy(fw, file)
					}
					continue
				}
			}
			writer.WriteField(pair[0], pair[1])
		} else {
			if marshaled, err := json.Marshal(unmarshalResult); err != nil {
				writer.WriteField(pair[0], pair[1])
			} else {
				writer.WriteField(part.Key, string(marshaled))
			}
		}
	}
	writer.Close()
	return HandleParseResult{
		ContentTypeHeader: writer.FormDataContentType(),
		Result:            &buf,
	}
}

func ParseUrlEncodedForm(data []HttpContentData) HandleParseResult {
	var formParts []string
	for _, part := range data {
		encodedKey := strings.ReplaceAll(strings.ReplaceAll(part.Key, " ", "+"), "=", "%3D")
		encodedValue := strings.ReplaceAll(strings.ReplaceAll(part.Value, " ", "+"), "=", "%3D")
		formParts = append(formParts, encodedKey+"="+encodedValue)
	}

	encoded := strings.Join(formParts, "&")

	return HandleParseResult{
		ContentTypeHeader: "application/x-www-form-urlencoded",
		Result:            encoded,
	}
}

func ParseHttpMethod(method string) string {
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

func ParseUrl(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return ""
	}
	return url
}

func ParseExecutionTimeInMilliseconds(start time.Time) float64 {
	executionTime := time.Since(start)
	ms := float64(executionTime.Nanoseconds()) / 1e6
	return ms
}

func DetectContentType(result string) string {
	trimmed := strings.TrimSpace(result)
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		var js json.RawMessage
		if json.Unmarshal([]byte(trimmed), &js) == nil {
			return "json"
		}
	}
	if strings.HasPrefix(trimmed, "<!DOCTYPE html") || strings.HasPrefix(trimmed, "<html") {
		return "html"
	}
	if strings.HasPrefix(trimmed, "<?xml") || (strings.HasPrefix(trimmed, "<") && strings.HasSuffix(trimmed, ">")) {
		return "xml"
	}
	return "text"
}

func GetFileByPath(path string) (*FileInfoData, error) {
	if len(strings.TrimSpace(path)) == 0 {
		return &FileInfoData{Name: "", PathIsValid: false}, os.ErrNotExist
	}
	if strings.ContainsAny(path, "\x00") {
		return &FileInfoData{Name: "", PathIsValid: false}, os.ErrInvalid
	}
	if !strings.ContainsAny(path, "/\\") {
		return &FileInfoData{Name: "", PathIsValid: false}, os.ErrInvalid
	}
	if !strings.ContainsAny(path, ".") {
		return &FileInfoData{Name: "", PathIsValid: false}, os.ErrInvalid
	}

	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return &FileInfoData{Name: "", PathIsValid: true}, os.ErrNotExist
	}
	
	return &FileInfoData{
		Name:        info.Name(),
		PathIsValid: true,
	}, nil
}
