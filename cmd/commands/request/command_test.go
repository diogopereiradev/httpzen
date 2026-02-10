package request_command

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/spf13/cobra"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
	http_utility "github.com/diogopereiradev/httpzen/internal/utils/http_utility"
)

type exitCalled struct{ code int }

func fakeExit(code int) {
	panic(exitCalled{code})
}

func Test_parseHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		want    http.Header
	}{
		{
			name:    "single header",
			headers: []string{"Content-Type: application/json"},
			want:    http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			name:    "multiple headers",
			headers: []string{"A: 1", "B: 2"},
			want:    http.Header{"A": []string{"1"}, "B": []string{"2"}},
		},
		{
			name:    "header with spaces",
			headers: []string{"  X-Test  :  value  "},
			want:    http.Header{"X-Test": []string{"value"}},
		},
		{
			name:    "invalid header (no colon)",
			headers: []string{"InvalidHeader"},
			want:    http.Header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseHeaders(tt.headers)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Init(t *testing.T) {
	var (
		calledBodyMenu    bool
		calledRunRequest  bool
		calledRequestMenu bool
		capturedOpts      request_module.RequestOptions
	)

	oldExit := Exit
	Exit = fakeExit
	defer func() { Exit = oldExit }()

	oldBodyMenu := BodyMenuNewFunc
	BodyMenuNewFunc = func(ro *request_module.RequestOptions, bodyPointer *[]http_utility.HttpContentData) {
		calledBodyMenu = true
	}
	defer func() { BodyMenuNewFunc = oldBodyMenu }()

	oldRunRequest := RunRequestFunc
	RunRequestFunc = func(opts request_module.RequestOptions) request_module.RequestResponse {
		calledRunRequest = true
		capturedOpts = opts
		return request_module.RequestResponse{}
	}
	defer func() { RunRequestFunc = oldRunRequest }()

	oldRequestMenu := RequestMenuNewFunc
	RequestMenuNewFunc = func(res *request_module.RequestResponse) {
		calledRequestMenu = true
	}
	defer func() { RequestMenuNewFunc = oldRequestMenu }()

	t.Run("calls help if not enough args", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)
		cmd.SetArgs([]string{})

		var helpCalled bool
		cmd.SetHelpFunc(func(*cobra.Command, []string) {
			helpCalled = true
		})

		cmd.Execute()
		if !helpCalled {
			t.Error("expected help to be called")
		}
	})

	t.Run("calls help if only method is provided", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)
		cmd.SetArgs([]string{"GET"})

		var helpCalled bool
		cmd.SetHelpFunc(func(*cobra.Command, []string) {
			helpCalled = true
		})

		cmd.Execute()
		if !helpCalled {
			t.Error("expected help to be called")
		}
	})

	t.Run("invalid HTTP method", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"INVALID", "http://test"})
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected exit to be called")
			}
		}()
		cmd.Execute()
	})

	t.Run("invalid URL", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"GET", "invalid-url"})
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected exit to be called")
			}
		}()
		cmd.Execute()
	})

	t.Run("body not allowed for GET", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"GET", "http://test"})
		cmd.Flags().Set("body", "true")
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected exit to be called")
			}
		}()
		cmd.Execute()
	})

	t.Run("body not allowed returns early when Exit is no-op", func(t *testing.T) {
		var exitCode int
		oldExitLocal := Exit
		Exit = func(code int) { exitCode = code }
		defer func() { Exit = oldExitLocal }()

		calledRunRequest = false
		calledRequestMenu = false
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)
		cmd.SetArgs([]string{"GET", "http://test", "foo=bar"})
		cmd.Execute()

		if exitCode != 1 {
			t.Errorf("expected exit code 1, got %d", exitCode)
		}
		if calledRunRequest || calledRequestMenu {
			t.Error("expected request to not run when body is not allowed")
		}
	})

	t.Run("body not allowed for HEAD", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"HEAD", "http://test"})
		cmd.Flags().Set("body", "true")
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected exit to be called")
			}
		}()
		cmd.Execute()
	})

	t.Run("valid request with body", func(t *testing.T) {
		calledBodyMenu = false
		calledRunRequest = false
		calledRequestMenu = false
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"POST", "http://test"})
		cmd.Flags().Set("body", "true")
		cmd.Execute()
		if !calledBodyMenu || !calledRunRequest || !calledRequestMenu {
			t.Error("expected all functions to be called")
		}
	})

	t.Run("valid request with inline body (does not open body menu)", func(t *testing.T) {
		calledBodyMenu = false
		calledRunRequest = false
		calledRequestMenu = false
		capturedOpts = request_module.RequestOptions{}

		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"POST", "http://test", "foo=\"bar\""})
		cmd.Flags().Set("insecure", "true")
		cmd.Execute()

		if calledBodyMenu {
			t.Error("expected body menu not to be called when inline body is provided")
		}
		if !calledRunRequest || !calledRequestMenu {
			t.Error("expected RunRequestFunc and RequestMenuNewFunc to be called")
		}
		if len(capturedOpts.Body) == 0 {
			t.Error("expected inline body to be set in request options")
		}
		if !capturedOpts.Insecure {
			t.Error("expected insecure flag to be propagated")
		}
	})

	t.Run("inline body not allowed for GET", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)
		cmd.SetArgs([]string{"GET", "http://test", "foo=bar"})
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected exit to be called")
			}
		}()
		cmd.Execute()
	})

	t.Run("valid request without body", func(t *testing.T) {
		calledRunRequest = false
		calledRequestMenu = false
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"GET", "http://test"})
		cmd.Execute()
		if !calledRunRequest || !calledRequestMenu {
			t.Error(
				"expected RunRequestFunc and RequestMenuNewFunc to be called",
			)
		}
	})

	t.Run("GET request with URL only (no method)", func(t *testing.T) {
		calledRunRequest = false
		calledRequestMenu = false
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"http://test"})
		cmd.Execute()
		if !calledRunRequest || !calledRequestMenu {
			t.Error(
				"expected RunRequestFunc and RequestMenuNewFunc to be called",
			)
		}
	})

	t.Run("GET request with https URL only (no method)", func(t *testing.T) {
		calledRunRequest = false
		calledRequestMenu = false
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)

		cmd.SetArgs([]string{"https://test"})
		cmd.Execute()
		if !calledRunRequest || !calledRequestMenu {
			t.Error(
				"expected RunRequestFunc and RequestMenuNewFunc to be called",
			)
		}
	})

	t.Run("GET request with shorthand port URL only (no method)", func(t *testing.T) {
		calledRunRequest = false
		calledRequestMenu = false
		capturedOpts = request_module.RequestOptions{}

		cmd := &cobra.Command{Use: "test"}
		Init(cmd)
		cmd.SetArgs([]string{":8080/api"})
		cmd.Execute()
		if !calledRunRequest || !calledRequestMenu {
			t.Error("expected RunRequestFunc and RequestMenuNewFunc to be called")
		}
		if capturedOpts.Url != "http://localhost:8080/api" {
			t.Errorf("expected url to be expanded, got %q", capturedOpts.Url)
		}
	})
}
