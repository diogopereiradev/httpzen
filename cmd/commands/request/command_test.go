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

	t.Run("valid request without body", func(t *testing.T) {
		calledRunRequest = false
		calledRequestMenu = false
		cmd := &cobra.Command{Use: "test"}
		Init(cmd)
		
		cmd.SetArgs([]string{"GET", "http://test"})
		cmd.Execute()
		if !calledRunRequest || !calledRequestMenu {
			t.Error("expected RunRequestFunc and RequestMenuNewFunc to be called")
		}
	})
}
