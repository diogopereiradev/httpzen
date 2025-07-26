package request_menu

import (
	"os"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	timed_message_component "github.com/diogopereiradev/httpzen/internal/components/timed_message"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

type Model struct {
	activeTab tab

	config   *config_module.Config
	response *request_module.RequestResponse

	clipboardTimedMessage *timed_message_component.TimedMessage

	isRefetching bool

	resultScrollOffset int
	resultLinesAmount  int

	infosScrollOffset int
	infosLinesAmount  int

	networkScrollOffset int
	networkLinesAmount  int

	reqHeadersScrollOffset int
	reqHeadersLinesAmount  int

	respHeadersScrollOffset int
	respHeadersLinesAmount  int
}

var Exit = os.Exit
var LoggerError = logger_module.Error
var TeaNewProgram = tea.NewProgram
var RunRequestFunc = request_module.RunRequest
var ClipboardWriteAll = clipboard.WriteAll
var TermClear = terminal_utility.Clear

var RunProgram = func(p *tea.Program) (tea.Model, error) {
	return p.Run()
}

func initialModel(res *request_module.RequestResponse, config *config_module.Config) Model {
	return Model{
		config:                config,
		activeTab:             tab_Result,
		response:              res,
		isRefetching:          false,
		clipboardTimedMessage: timed_message_component.NewTimedMessage(),
	}
}

func New(res *request_module.RequestResponse) {
	config := config_module.GetConfig()
	model := initialModel(res, &config)

	p := TeaNewProgram(&model)
	TermClear()

	if _, err := RunProgram(p); err != nil {
		LoggerError("Error on rendering the program: " + err.Error())
		Exit(1)
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var content string

	if m.isRefetching {
		content = refetch_Render(m)
		return content
	}

	if !m.config.HideLogomark {
		content += lipgloss.NewStyle().Foreground(theme.Primary).Render(logoascii.GetLogo(".request")) + "\n"
	}

	content += tab_Render(m)
	content += "\n\n"

	switch m.activeTab {
	case tab_Result:
		content += result_viewport_Render(m)
	case tab_RequestInfos:
		content += basic_infos_Render_Paged(m)
	case tab_NetworkInfos:
		content += network_infos_Render_Paged(m)
	case tab_RequestHeaders:
		content += request_headers_Render_Paged(m)
	case tab_ResponseHeaders:
		content += response_headers_Render_Paged(m)
	}
	content += navigation_options_Render()

	if m.clipboardTimedMessage != nil && m.clipboardTimedMessage.Visible {
		dialogMsg := m.clipboardTimedMessage.Render()
		if dialogMsg != "" {
			content += dialogMsg
		}
	}
	return content
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Events
	switch ev := msg.(type) {
	case RefetchEvent:
		model := initialModel(&ev.Response, m.config)
		model.activeTab = m.activeTab
		model.isRefetching = false
		m.isRefetching = false
		m = &model
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Shortcuts
		switch keyMsg.Type {
		case tea.KeyRunes:
			if keyMsg.String() == "q" {
				return m, tea.Quit
			}
			if keyMsg.String() == "r" {
				m.isRefetching = true
				return m, func() tea.Msg {
					res := RunRequestFunc(request_module.RequestOptions{
						Url:     m.response.Request.Url,
						Headers: m.response.Request.Headers,
						Method:  m.response.Request.Method,
						Timeout: m.response.Request.Timeout,
						Body:    m.response.Request.Body,
					})
					return RefetchEvent{Response: res}
				}
			}
			if keyMsg.String() == "c" {
				err := ClipboardWriteAll(m.response.Result)
				if err != nil {
					panic(err)
				} else {
					return m, m.clipboardTimedMessage.Show("Request response copied", 1*time.Second)
				}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyRight:
			tab_MoveRight(m)
		case tea.KeyLeft:
			tab_MoveLeft(m)

		case tea.KeyUp:
			switch m.activeTab {
			case tab_Result:
				result_viewport_ScrollUp(m)
			case tab_RequestInfos:
				basic_infos_ScrollUp(m)
			case tab_NetworkInfos:
				network_infos_ScrollUp(m)
			case tab_RequestHeaders:
				request_headers_ScrollUp(m)
			case tab_ResponseHeaders:
				response_headers_ScrollUp(m)
			}
		case tea.KeyDown:
			switch m.activeTab {
			case tab_Result:
				result_viewport_ScrollDown(m)
			case tab_RequestInfos:
				basic_infos_ScrollDown(m)
			case tab_NetworkInfos:
				network_infos_ScrollDown(m)
			case tab_RequestHeaders:
				request_headers_ScrollDown(m)
			case tab_ResponseHeaders:
				response_headers_ScrollDown(m)
			}
		case tea.KeyPgUp:
			switch m.activeTab {
			case tab_Result:
				result_viewport_ScrollPgUp(m)
			case tab_RequestInfos:
				basic_infos_ScrollPgUp(m)
			case tab_NetworkInfos:
				network_infos_ScrollPgUp(m)
			case tab_RequestHeaders:
				request_headers_ScrollPgUp(m)
			case tab_ResponseHeaders:
				response_headers_ScrollPgUp(m)
			}
		case tea.KeyPgDown:
			switch m.activeTab {
			case tab_Result:
				result_viewport_ScrollPgDown(m)
			case tab_RequestInfos:
				basic_infos_ScrollPgDown(m)
			case tab_NetworkInfos:
				network_infos_ScrollPgDown(m)
			case tab_RequestHeaders:
				request_headers_ScrollPgDown(m)
			case tab_ResponseHeaders:
				response_headers_ScrollPgDown(m)
			}
		}
	}
	return m, nil
}
