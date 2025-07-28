package http_utility

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseApplicationJson(t *testing.T) {
	valid := HttpContentData{ContentType: "application/json", Value: `{"foo":"bar"}`}

	res := ParseApplicationJson(valid)
	if res.ContentTypeHeader != "application/json" {
		t.Errorf("expected content type header")
	}

	if res.Result == nil {
		t.Errorf("expected result to be not nil")
	}

	invalid := HttpContentData{ContentType: "application/json", Value: `{"foo":}`}
	res = ParseApplicationJson(invalid)
	if res.Result != nil {
		t.Errorf("expected result to be nil on error")
	}
}

func TestParseMultipartFormData(t *testing.T) {
	file, err := os.CreateTemp("", "testfile*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString("testdata")
	file.Close()

	data := []HttpContentData{
		{Key: "field1", Value: "value1"},
		{Key: "file1", Value: file.Name()},
	}

	res := ParseMultipartFormData(data)
	if !strings.HasPrefix(res.ContentTypeHeader, "multipart/form-data;") {
		t.Errorf("expected multipart content type")
	}

	buf, ok := res.Result.(*bytes.Buffer)
	if !ok || buf.Len() == 0 {
		t.Errorf("expected buffer result")
	}
}

func TestParseMultipartFormData_FileOpenError(t *testing.T) {
	data := []HttpContentData{
		{Key: "file1", Value: "/path/that/does/not/exist.txt"},
	}

	res := ParseMultipartFormData(data)
	if !strings.HasPrefix(res.ContentTypeHeader, "multipart/form-data;") {
		t.Errorf("expected multipart content type")
	}
}

func TestParseMultipartFormData_CreateFormFileError(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdir*")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)
	data := []HttpContentData{
		{Key: "file1", Value: dir},
	}

	res := ParseMultipartFormData(data)
	if !strings.HasPrefix(res.ContentTypeHeader, "multipart/form-data;") {
		t.Errorf("expected multipart content type")
	}
}

func TestParseMultipartFormData_MarshalError(t *testing.T) {
	data := []HttpContentData{
		{Key: "field1", Value: `{"foo":}`},
	}

	res := ParseMultipartFormData(data)
	if !strings.HasPrefix(res.ContentTypeHeader, "multipart/form-data;") {
		t.Errorf("expected multipart content type")
	}
}

func TestParseMultipartFormData_ValidPath(t *testing.T) {
	data := []HttpContentData{
		{Key: "file1", Value: "./util.go"},
	}

	res := ParseMultipartFormData(data)
	if !strings.HasPrefix(res.ContentTypeHeader, "multipart/form-data;") {
		t.Errorf("expected multipart content type")
	}

	buf, ok := res.Result.(*bytes.Buffer)
	if !ok || buf.Len() == 0 {
		t.Errorf("expected buffer result")
	}
}

func TestParseMultipartFormData_ValidUnmarshal(t *testing.T) {
	data := []HttpContentData{
		{Key: "field1", Value: `{"foo":"bar"}`},
	}

	res := ParseMultipartFormData(data)
	if !strings.HasPrefix(res.ContentTypeHeader, "multipart/form-data;") {
		t.Errorf("expected multipart content type")
	}

	buf, ok := res.Result.(*bytes.Buffer)
	if !ok || buf.Len() == 0 {
		t.Errorf("expected buffer result")
	}
}

func TestParseUrlEncodedForm(t *testing.T) {
	data := []HttpContentData{
		{Key: "foo bar", Value: "baz=qux"},
	}

	res := ParseUrlEncodedForm(data)
	if res.ContentTypeHeader != "application/x-www-form-urlencoded" {
		t.Errorf("expected urlencoded content type")
	}

	if !strings.Contains(res.Result.(string), "foo+bar=baz%3Dqux") {
		t.Errorf("expected encoded result")
	}
}

func TestParseHttpMethod(t *testing.T) {
	cases := map[string]string{
		"get":    "GET",
		"post":   "POST",
		"put":    "PUT",
		"delete": "DELETE",
		"patch":  "PATCH",
		"head":   "HEAD",
		"other":  "",
	}

	for in, want := range cases {
		if got := ParseHttpMethod(in); got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	}
}

func TestParseUrl(t *testing.T) {
	if ParseUrl("http://foo") != "http://foo" {
		t.Errorf("expected valid url")
	}

	if ParseUrl("https://foo") != "https://foo" {
		t.Errorf("expected valid url")
	}

	if ParseUrl("ftp://foo") != "" {
		t.Errorf("expected empty for invalid url")
	}
}

func TestParseExecutionTimeInMilliseconds(t *testing.T) {
	start := time.Now()
	time.Sleep(10 * time.Millisecond)

	ms := ParseExecutionTimeInMilliseconds(start)
	if ms < 0 {
		t.Errorf("expected positive ms")
	}
}

func TestDetectContentType(t *testing.T) {
	if DetectContentType("{\"foo\":1}") != "json" {
		t.Errorf("expected json")
	}

	if DetectContentType("<html>") != "html" {
		t.Errorf("expected html")
	}

	if DetectContentType("<?xml version='1.0'?>") != "xml" {
		t.Errorf("expected xml")
	}

	if DetectContentType("plain text") != "text" {
		t.Errorf("expected text")
	}
}

func TestGetFileByPath(t *testing.T) {
	_, err := GetFileByPath("")
	if err == nil {
		t.Errorf("expected error for empty path")
	}

	_, err = GetFileByPath("\x00badpath")
	if err == nil {
		t.Errorf("expected error for null byte")
	}

	_, err = GetFileByPath("no_slash.txt")
	if err == nil {
		t.Errorf("expected error for missing slash")
	}

	_, err = GetFileByPath("/tmp/nodot")
	if err == nil {
		t.Errorf("expected error for missing dot")
	}

	file, err := os.CreateTemp("/tmp", "testfile*.txt")
	if err != nil {
		t.Fatal(err)
	}

	file.Close()

	defer os.Remove(file.Name())
	info, err := GetFileByPath(file.Name())
	if err != nil {
		t.Errorf("expected valid file")
	}

	if !info.PathIsValid {
		t.Errorf("expected valid path")
	}
}
