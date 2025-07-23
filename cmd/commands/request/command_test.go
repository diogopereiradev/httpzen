package request_command

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
)

var (
	mockExitCalled      bool
	mockExitCode        int
	mockRunRequestInput interface{}
	mockRunRequestOut   interface{}
	mockMenuCalled      bool
)

func resetMocks() {
	mockExitCalled = false
	mockExitCode = 0
	mockRunRequestInput = nil
	mockRunRequestOut = request_module.RequestResponse{}
	mockMenuCalled = false
}

func mockExit(code int) {
	mockExitCalled = true
	mockExitCode = code
	panic("exit")
}

func mockRunRequest(opt request_module.RequestOptions) request_module.RequestResponse {
	mockRunRequestInput = opt
	return mockRunRequestOut.(request_module.RequestResponse)
}

func mockMenu(res *request_module.RequestResponse) {
	mockMenuCalled = true
}

func TestInit_HelpCalledOnMissingArgs(t *testing.T) {
	resetMocks()
	cmd := &cobra.Command{Use: "test"}
	helpCalled := false
	cmd.SetHelpFunc(func(*cobra.Command, []string) {
		helpCalled = true
	})
	Init(cmd)
	cmd.SetArgs([]string{})
	cmd.Execute()
	assert.True(t, helpCalled)
}

func TestInit_InvalidMethod(t *testing.T) {
	resetMocks()
	Exit = mockExit
	defer func() { Exit = exitOrig }()
	cmd := &cobra.Command{Use: "test"}
	Init(cmd)
	cmd.SetArgs([]string{"INVALID", "http://url"})
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "exit", r)
		}
		assert.True(t, mockExitCalled)
		assert.Equal(t, 1, mockExitCode)
	}()
	cmd.Execute()
}

func TestInit_InvalidUrl(t *testing.T) {
	resetMocks()
	Exit = mockExit
	defer func() { Exit = exitOrig }()
	cmd := &cobra.Command{Use: "test"}
	Init(cmd)
	cmd.SetArgs([]string{"GET", "badurl"})
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "exit", r)
		}
		assert.True(t, mockExitCalled)
		assert.Equal(t, 1, mockExitCode)
	}()
	cmd.Execute()
}

func TestInit_ValidRequest_AllFlags(t *testing.T) {
	resetMocks()
	Exit = mockExit
	RunRequestFunc = mockRunRequest
	RequestMenuNewFunc = mockMenu
	defer func() {
		Exit = exitOrig
		RunRequestFunc = runRequestOrig
		RequestMenuNewFunc = menuOrig
	}()
	cmd := &cobra.Command{Use: "test"}
	Init(cmd)
	cmd.SetArgs([]string{"GET", "http://url"})
	cmd.Flags().Set("header", "X-Test:foo")
	cmd.Flags().Set("json", "true")
	cmd.Flags().Set("raw", "true")
	cmd.Flags().Set("form", "true")
	cmd.Flags().Set("multipart", "true")
	cmd.Flags().Set("headers", "true")
	cmd.Flags().Set("body", "true")
	cmd.Flags().Set("meta", "true")
	cmd.Execute()
	opt, ok := mockRunRequestInput.(request_module.RequestOptions)
	assert.True(t, ok)
	assert.Equal(t, "http://url", opt.Url)
	assert.True(t, mockMenuCalled)
}

func Test_parseHeaders(t *testing.T) {
	h := []string{"A:1", "B:2:3", "C : 4"}
	head := parseHeaders(h)
	assert.Equal(t, "1", head.Get("A"))
	assert.Equal(t, "2:3", head.Get("B"))
	assert.Equal(t, "4", head.Get("C"))
	assert.Equal(t, 3, len(head))
}

var exitOrig = Exit
var runRequestOrig = RunRequestFunc
var menuOrig = RequestMenuNewFunc
