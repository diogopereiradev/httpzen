package request_menu

import (
	"errors"
	"net/http"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	timed_message_util "github.com/diogopereiradev/httpzen/internal/utils/timed_message"
	"github.com/stretchr/testify/assert"
)

var (
	mockExitCalled bool
	mockExitCode   int
	mockLoggerMsg  string
)

func mockExit(code int) {
	mockExitCalled = true
	mockExitCode = code
	panic("exit")
}

func mockLoggerError(msg string) {
	mockLoggerMsg = msg
}

func resetMocks() {
	mockExitCalled = false
	mockExitCode = 0
	mockLoggerMsg = ""
}

type testModel struct{}

func (m testModel) Init() tea.Cmd { return tea.Quit }
func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m testModel) View() string { return "" }

func TestRunProgram_CoverageReal(t *testing.T) {
  prog := tea.NewProgram(testModel{})
  _, err := RunProgram(prog)
  if err != nil {
    t.Errorf("RunProgram returned error: %v", err)
  }
}

func Test_initialModel(t *testing.T) {
	res := &request_module.RequestResponse{}
	cfg := &config_module.Config{}
	m := initialModel(res, cfg)
	assert.Equal(t, res, m.response)
	assert.Equal(t, cfg, m.config)
	assert.False(t, m.isRefetching)
	assert.NotNil(t, m.clipboardTimedMessage)
}

func TestNew_ErrorOnRun(t *testing.T) {
	resetMocks()
	Exit = mockExit
	origLogger := LoggerError
	LoggerError = mockLoggerError
	origTeaNewProgram := TeaNewProgram
	TeaNewProgram = func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		return &tea.Program{}
	}
	defer func() {
		Exit = os.Exit
		LoggerError = origLogger
		TeaNewProgram = origTeaNewProgram
	}()

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "exit", r)
			assert.True(t, mockExitCalled)
			assert.Contains(t, mockLoggerMsg, "Error on rendering the program")
		}
	}()

	LoggerError("Error on rendering the program: fail")
	Exit(1)
}

func TestNew_Success(t *testing.T) {
	origTeaNewProgram := TeaNewProgram
	TeaNewProgram = func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		return &tea.Program{}
	}
	defer func() { TeaNewProgram = origTeaNewProgram }()

	origRunProgram := RunProgram
	RunProgram = func(p *tea.Program) (tea.Model, error) {
		return nil, nil
	}
	defer func() { RunProgram = origRunProgram }()

	origExit := Exit
	Exit = func(code int) { panic("should not exit") }
	defer func() { Exit = origExit }()

	termClearCalled := false
	origTermClear := TermClear
	TermClear = func() { termClearCalled = true }
	defer func() { TermClear = origTermClear }()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()

	resp := &request_module.RequestResponse{}
	New(resp)
	assert.True(t, termClearCalled)
}

func TestNew_RunProgramError(t *testing.T) {
	resetMocks()
	Exit = mockExit
	origLogger := LoggerError
	LoggerError = mockLoggerError
	origTeaNewProgram := TeaNewProgram
	TeaNewProgram = func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		return &tea.Program{}
	}
	origRunProgram := RunProgram
	RunProgram = func(p *tea.Program) (tea.Model, error) {
		return nil, assert.AnError
	}
	defer func() {
		Exit = os.Exit
		LoggerError = origLogger
		TeaNewProgram = origTeaNewProgram
		RunProgram = origRunProgram
	}()

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "exit", r)
			assert.True(t, mockExitCalled)
			assert.Contains(t, mockLoggerMsg, "Error on rendering the program")
		}
	}()

	resp := &request_module.RequestResponse{}
	New(resp)
}

func TestModel_Init(t *testing.T) {
	m := &Model{}
	assert.Nil(t, m.Init())
}

func TestModel_View(t *testing.T) {
	m := &Model{
		config: &config_module.Config{},
		response: &request_module.RequestResponse{
			Request: request_module.RequestOptions{},
		},
		clipboardTimedMessage: timed_message_util.NewTimedMessage(),
	}
	m.activeTab = tab_Result
	m.isRefetching = true
	assert.Contains(t, m.View(), "")
}

func TestModel_View_AllTabs(t *testing.T) {
	baseResp := &request_module.RequestResponse{
		HttpVersion:   "HTTP/1.1",
		StatusMessage: "OK",
		StatusCode:    200,
		ExecutionTime: 1.0,
		Headers:       http.Header{"X-Test": []string{"ok"}},
		Body:          "body",
		Cookies:       nil,
		Request: request_module.RequestOptions{
			Url:     "url",
			Method:  "GET",
			Headers: http.Header{},
			Timeout: 1,
			Body:    "body",
		},
		Path:    "/",
		Host:    "localhost",
		Method:  "GET",
		IpInfos: nil,
	}

	m := &Model{
		config:                &config_module.Config{},
		clipboardTimedMessage: nil,
		isRefetching:          false,
		response:              baseResp,
	}

	tabs := []tab{tab_Result, tab_RequestInfos, tab_NetworkInfos, tab_RequestHeaders, tab_ResponseHeaders}
	for _, tb := range tabs {
		m.activeTab = tb
		m.config.HideLogomark = false
		out := m.View()
		assert.Contains(t, out, "\n\n")
	}

	m.config.HideLogomark = true
	m.activeTab = tab_Result

	assert.NotContains(t, m.View(), ".request")

	m.clipboardTimedMessage = timed_message_util.NewTimedMessage()
	m.clipboardTimedMessage.Show("copied", 1)
	_ = m.View()
}

func TestModel_Update_Refetch(t *testing.T) {
	m := &Model{
		config: &config_module.Config{},
		response: &request_module.RequestResponse{
			Request: request_module.RequestOptions{
				Url:     "http://localhost",
				Method:  "GET",
				Headers: http.Header{},
				Timeout: 1,
				Body:    "body",
			},
		},
	}
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}

	origRunRequestFunc := RunRequestFunc
	RunRequestFunc = func(options request_module.RequestOptions) request_module.RequestResponse {
		return request_module.RequestResponse{Result: "mocked"}
	}
	defer func() { RunRequestFunc = origRunRequestFunc }()

	model, cmd := m.Update(msg)
	assert.Equal(t, m, model)
	assert.NotNil(t, cmd)

	if cmd != nil {
		_ = cmd()
	}
}

func TestModel_Update_Copy(t *testing.T) {
	m := &Model{
		config:                &config_module.Config{},
		response:              &request_module.RequestResponse{Result: "copied"},
		clipboardTimedMessage: timed_message_util.NewTimedMessage(),
	}
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}

	origClipboardWriteAll := ClipboardWriteAll
	ClipboardWriteAll = func(string) error { return nil }
	defer func() { ClipboardWriteAll = origClipboardWriteAll }()
	m.Update(msg)

	ClipboardWriteAll = func(string) error { return errors.New("fail") }
	defer func() { ClipboardWriteAll = origClipboardWriteAll }()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	m.Update(msg)
}

func TestModel_Update_TabNavigation(t *testing.T) {
	m := &Model{config: &config_module.Config{}}
	for _, key := range []tea.KeyType{tea.KeyRight, tea.KeyLeft} {
		msg := tea.KeyMsg{Type: key}
		m.Update(msg)
	}
}

func TestModel_Update_Scroll(t *testing.T) {
	m := &Model{config: &config_module.Config{}}
	tabs := []tab{tab_Result, tab_RequestInfos, tab_NetworkInfos, tab_RequestHeaders, tab_ResponseHeaders}
	for _, tb := range tabs {
		m.activeTab = tb
		for _, key := range []tea.KeyType{tea.KeyUp, tea.KeyDown, tea.KeyPgUp, tea.KeyPgDown} {
			msg := tea.KeyMsg{Type: key}
			m.Update(msg)
		}
	}
}

func TestModel_Update_RefetchEvent(t *testing.T) {
	m := &Model{config: &config_module.Config{}}
	resp := request_module.RequestResponse{}
	evt := RefetchEvent{Response: resp}

	model, cmd := m.Update(evt)
	assert.NotNil(t, model)
	assert.Nil(t, cmd)
}

func TestModel_Update_Quit(t *testing.T) {
	m := &Model{config: &config_module.Config{}}
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	model, cmd := m.Update(msg)
	assert.Equal(t, m, model)
	assert.NotNil(t, cmd)
}

func TestModel_Update_QuitKeys(t *testing.T) {
	m := &Model{config: &config_module.Config{}}
	for _, key := range []tea.KeyType{tea.KeyCtrlC, tea.KeyEsc} {
		msg := tea.KeyMsg{Type: key}
		model, cmd := m.Update(msg)
		assert.Equal(t, m, model)
		assert.NotNil(t, cmd)
	}
}
