package request_menu

import (
	"os"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/term_clear"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
	timed_message_util "github.com/diogopereiradev/httpzen/internal/utils/timed_message"
)

type Model struct {
	config    *config_module.Config
	activeTab tab
	response  *request_module.RequestResponse
	clipbloardTimedMessage *timed_message_util.TimedMessage

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

func initialModel(res *request_module.RequestResponse, config *config_module.Config) Model {
	return Model{
		config:                 config,
		activeTab:              tab_Result,
		response:               res,
		clipbloardTimedMessage: timed_message_util.NewTimedMessage(),
	}
}

func New(res *request_module.RequestResponse) {
	config := config_module.GetConfig()
	model := initialModel(res, &config)

	p := tea.NewProgram(&model)

	term_clear.Clear()

	if _, err := p.Run(); err != nil {
		logger_module.Error("Error on rendering the program: " + err.Error())
		Exit(1)
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var content string

	if !m.config.HideLogomark {
		content += lipgloss.NewStyle().Foreground(theme.Primary).Render(logoascii.GetLogo(".request")) + "\n"
	}

	content += tab_Render(m)
	content += "\n\n"

	switch m.activeTab {
	case tab_Result: content += result_viewport_Render(m)
	case tab_RequestInfos: content += basic_infos_Render_Paged(m)
	case tab_NetworkInfos: content += network_infos_Render_Paged(m)
	case tab_RequestHeaders: content += request_headers_Render_Paged(m)
	case tab_ResponseHeaders: content += response_headers_Render_Paged(m)
	}
	content += navigation_options_Render()

	if m.clipbloardTimedMessage != nil && m.clipbloardTimedMessage.Visible {
		dialogMsg := m.clipbloardTimedMessage.Render()
		if dialogMsg != "" {
			content += dialogMsg
		}
	}

	return content
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		  case tea.KeyRunes:
		  	if keyMsg.String() == "q" {
		  		return m, tea.Quit
		  	}
		  	if keyMsg.String() == "c" {
		  		err := clipboard.WriteAll(m.response.Result)
		  		if err != nil {
		  			panic(err)
		  		} else {
		  			return m, m.clipbloardTimedMessage.Show("Request response copied", 1 * time.Second)
		  		}
		  	}
  
		  case tea.KeyCtrlC, tea.KeyEsc:
		  	return m, tea.Quit
  
		  case tea.KeyRight: tab_MoveRight(m)
		  case tea.KeyLeft: tab_MoveLeft(m)
  
		  case tea.KeyUp:
		  	switch m.activeTab {
		  	  case tab_Result: result_viewport_ScrollUp(m)
		  	  case tab_RequestInfos: basic_infos_ScrollUp(m)
		  	  case tab_NetworkInfos: network_infos_ScrollUp(m)
		  	  case tab_RequestHeaders: request_headers_ScrollUp(m)
		  	  case tab_ResponseHeaders: response_headers_ScrollUp(m)
		  	}
		  case tea.KeyDown:
		  	switch m.activeTab {
		  	  case tab_Result: result_viewport_ScrollDown(m)
		  	  case tab_RequestInfos: basic_infos_ScrollDown(m)
		  	  case tab_NetworkInfos: network_infos_ScrollDown(m)
		  	  case tab_RequestHeaders: request_headers_ScrollDown(m)
		  	  case tab_ResponseHeaders: response_headers_ScrollDown(m)
		  	}
		  case tea.KeyPgUp:
		  	switch m.activeTab {
		  	  case tab_Result: result_viewport_ScrollPgUp(m)
		  	  case tab_RequestInfos: basic_infos_ScrollPgUp(m)
		  	  case tab_NetworkInfos: network_infos_ScrollPgUp(m)
		  	  case tab_RequestHeaders: request_headers_ScrollPgUp(m)
		  	  case tab_ResponseHeaders: response_headers_ScrollPgUp(m)
		  	}
		  case tea.KeyPgDown:
		  	switch m.activeTab {
		  	  case tab_Result: result_viewport_ScrollPgDown(m)
		  	  case tab_RequestInfos: basic_infos_ScrollPgDown(m)
		  	  case tab_NetworkInfos: network_infos_ScrollPgDown(m)
		  	  case tab_RequestHeaders: request_headers_ScrollPgDown(m)
		  	  case tab_ResponseHeaders: response_headers_ScrollPgDown(m)
		  	}
		}
	}
	return m, nil
}
