package benchmark_menu

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/diogopereiradev/httpzen/internal/utils/theme"
)

func metrics_Render(m *BenchmarkModel) string {
	var content string

	fieldStyle := lipgloss.NewStyle().Foreground(theme.Secondary)
	greyStyle := lipgloss.NewStyle().Foreground(theme.DarkenText)

	content += fieldStyle.Render("Threads: ") + i32toa(m.metrics.ExecutedThreads) + "\n"
	content += fieldStyle.Render("Total Requests: ") + itoa(m.metrics.TotalRequests) + "\n"
	content += fieldStyle.Render("Total Success: ") + itoa(m.metrics.TotalSuccess) + "\n"
	content += fieldStyle.Render("Total Errors: ") + itoa(m.metrics.TotalErrors) + "\n"
	content += fieldStyle.Render("Total data sent: ") + bytesToMb(m.metrics.TotalBytesSent) + "\n"
	content += fieldStyle.Render("Total data received: ") + bytesToMb(m.metrics.TotalBytesReceived) + "\n"
	content += fieldStyle.Render("Total Duration: ") + itoa(m.metrics.Duration) + "/" + greyStyle.Render(itoa(m.metrics.TotalDuration)) + " seconds" + "\n"
	content += fieldStyle.Render("Min Latency: ") + latencyFormat(m.metrics.RequestsMinLatency) + "\n"
	content += fieldStyle.Render("Max Latency: ") + latencyFormat(m.metrics.RequestsMaxLatency) + "\n"
	content += fieldStyle.Render("Requests per second: ") + itoa(m.metrics.RequestsPerSecond) + "\n"

	return content
}

func i32toa(i int32) string {
	return fmt.Sprintf("%d", i)
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func bytesToMb(b int) string {
	return fmt.Sprintf("%.2f MB", float64(b)/1024/1024)
}

func latencyFormat(d float64) string {
	return fmt.Sprintf("%.2f ms", d)
}