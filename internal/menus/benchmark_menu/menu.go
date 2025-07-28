package benchmark_menu

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	benchmark_module "github.com/diogopereiradev/httpzen/internal/benchmark"
	config_module "github.com/diogopereiradev/httpzen/internal/config"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	logoascii "github.com/diogopereiradev/httpzen/internal/utils/logo_ascii"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

var New = newComponent

type BenchmarkModel struct {
	ThreadsAmount int
	Duration      int
	Request       request_module.RequestOptions

	metrics      *benchmark_module.Metrics
	config       config_module.Config
	benchmarking bool
	done         bool
}

type benchmarkProgressMsg struct {
	metrics *benchmark_module.Metrics
}

type benchmarkResultMsg struct {
	metrics benchmark_module.Metrics
}

func newComponent(model *benchmark_module.BenchmarkOptions) {
	m := &BenchmarkModel{
		Request:       model.Request,
		ThreadsAmount: model.ThreadsAmount,
		Duration:      model.Duration,
		metrics:       &benchmark_module.Metrics{},
		config:        config_module.GetConfig(),
	}

	terminal_utility.Clear()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

func (m *BenchmarkModel) Init() tea.Cmd {
	return m.runBenchmarkRealtime()
}

func (m *BenchmarkModel) View() string {
	var content string

	titleStyle := lipgloss.NewStyle().Foreground(theme.Primary)
	greyStyle := lipgloss.NewStyle().Foreground(theme.DarkenText)
	borderStyle := lipgloss.
		NewStyle().
		Width(terminal_utility.GetTerminalWidth(60)-2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		Padding(0, 2, 1, 2)

	labeledStyle := lipgloss.NewStyle().
		Foreground(theme.LightText).
		Padding(0, 2).
		Bold(true)

	if !m.config.HideLogomark {
		content += titleStyle.Render(logoascii.GetLogo(".benchmark")) + "\n"
	} else {
		borderStyle = borderStyle.Padding(1, 2, 1, 2)
	}

	if m.done {
		content += metrics_Render(m)
		content += "\n" + labeledStyle.Background(theme.Success).Render("Benchmark completed!")
	}

	if m.benchmarking {
		content += metrics_Render(m)
		content += "\n" + labeledStyle.Background(theme.Primary).Render("Benchmarking...")
	}

	content += greyStyle.Render("\n\nYou can press 'q' or 'ctrl+c' to cancel/quit.")
	return borderStyle.Render(content) + "\n"
}

func (m *BenchmarkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyCtrlC: return m, tea.Quit
		case tea.KeyRunes:
			if keyMsg.String() == "q" { return m, tea.Quit }
		}
	}

	switch msg := msg.(type) {
	case benchmarkResultMsg:
		m.metrics = &msg.metrics
		m.benchmarking = false
		m.done = true
		return m, nil
	case benchmarkProgressMsg:
		m.metrics = msg.metrics
		m.metrics.TotalDuration = m.Duration
		return m, tea.Tick(1000*time.Millisecond, func(_ time.Time) tea.Msg {
			if m.metrics.Duration >= m.Duration {
				return benchmarkResultMsg{metrics: *m.metrics}
			}
			return benchmarkProgressMsg{metrics: msg.metrics}
		})
	}
	return m, nil
}

func (m *BenchmarkModel) runBenchmarkRealtime() tea.Cmd {
	m.benchmarking = true
	return func() tea.Msg {
		metrics := &benchmark_module.Metrics{}
		done := make(chan struct{})

		go func() {
			defer close(done)
			options := benchmark_module.BenchmarkOptions{
				Request:       m.Request,
				ThreadsAmount: m.ThreadsAmount,
				Duration:      m.Duration,
			}
			benchmark_module.RunBenchmark(options, metrics)
		}()

		realtimeTicker := time.NewTicker(1000 * time.Millisecond)
		defer realtimeTicker.Stop()

		for {
			select {
			case <-done:
				return benchmarkResultMsg{metrics: *metrics}
			case <-realtimeTicker.C:
				return benchmarkProgressMsg{metrics: metrics}
			}
		}
	}
}
