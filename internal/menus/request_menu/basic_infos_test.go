package request_menu

import (
	"strings"
	"testing"

	config_module "github.com/diogopereiradev/httpzen/internal/config"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
)

func MinTestHelper(a, b int) int {
	return min(a, b)
}

func min(a, b int) int {
	if a < 0 {
		return 0
	}
	if b < 0 {
		return 0
	}
	if a < b {
		return a
	}
	return b
}

func makeModel() *Model {
	return &Model{
		config: &config_module.Config{SlowResponseThreshold: 50},
		response: &request_module.RequestResponse{
			ExecutionTime: 100,
			HttpVersion:   "HTTP/1.1",
			Method:        "GET",
			StatusMessage: "200 OK",
			Result:        strings.Repeat("a", 200),
			Body:          []request_module.RequestBody{{Key: "body", Value: "request body"}},
		},
	}
}

func Test_basic_infos_Render(t *testing.T) {
	m := makeModel()
	out := basic_infos_Render(m)
	if !strings.Contains(out, "HTTP/1.1 GET 200 OK") {
		t.Error("Should contain status in the first line")
	}
	if !strings.Contains(out, "Response Time:") {
		t.Error("Should contain response time")
	}
	if !strings.Contains(out, "Response Size:") {
		t.Error("Should contain response size")
	}
	if !strings.Contains(out, "request body") {
		t.Error("Should contain request body if present")
	}
}

func Test_basic_infos_Render_Paged(t *testing.T) {
	m := makeModel()
	m.response.Result = strings.Repeat("line\n", 30)

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 5, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()

	out := basic_infos_Render_Paged(m)
	t.Log("EXIT OF RENDER_PAGED:\n" + out)

	if m.infosLinesAmount == 0 {
		t.Error("infosLinesAmount should be set")
	}

	if !strings.Contains(out, "Use ↑/↓ or PgUp/PgDown to scroll.") {
		t.Error("Should show scroll instruction if there are more lines than the terminal height")
	}
}

func Test_basic_infos_ScrollUp(t *testing.T) {
	m := makeModel()
	m.infosScrollOffset = 2

	basic_infos_ScrollUp(m)
	if m.infosScrollOffset != 1 {
		t.Errorf("Expected 1, got %d", m.infosScrollOffset)
	}

	basic_infos_ScrollUp(m)
	basic_infos_ScrollUp(m)
	if m.infosScrollOffset != 0 {
		t.Errorf("Expected 0, got %d", m.infosScrollOffset)
	}
}

func Test_basic_infos_ScrollDown(t *testing.T) {
	m := makeModel()
	m.infosLinesAmount = 100

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 10, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()

	m.infosScrollOffset = 0

	basic_infos_ScrollDown(m)
	if m.infosScrollOffset != 1 {
		t.Errorf("Expected 1, got %d", m.infosScrollOffset)
	}
}

func Test_basic_infos_ScrollDown_NoLines(t *testing.T) {
	m := makeModel()
	m.infosLinesAmount = 0

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 10, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()

	m.infosScrollOffset = 0

	basic_infos_ScrollDown(m)
	if m.infosScrollOffset != 0 {
		t.Errorf("Should not advance scroll if infosLinesAmount is 0")
	}
}

func Test_basic_infos_ScrollDown_AtEnd(t *testing.T) {
	m := makeModel()
	m.infosLinesAmount = 15

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 10, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()
	term_size.GetHeightFunc = func() (int, error) { return 20, nil }

	m.infosScrollOffset = 11

	basic_infos_ScrollDown(m)
	if m.infosScrollOffset != 11 {
		t.Errorf("Should not advance scroll if already at end")
	}
}

func Test_basic_infos_ScrollDown_TooFewLinesForScroll(t *testing.T) {
	m := makeModel()
	m.infosLinesAmount = 5

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 50, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()
	m.infosScrollOffset = 0

	basic_infos_ScrollDown(m)
	if m.infosScrollOffset != 0 {
		t.Errorf("Should not advance scroll if not enough lines to scroll")
	}
}

func Test_basic_infos_ScrollPgUp(t *testing.T) {
	m := makeModel()
	m.infosScrollOffset = 6

	basic_infos_ScrollPgUp(m)
	if m.infosScrollOffset != 1 {
		t.Errorf("Expected 1, got %d", m.infosScrollOffset)
	}

	basic_infos_ScrollPgUp(m)
	if m.infosScrollOffset != 0 {
		t.Errorf("Expected 0, got %d", m.infosScrollOffset)
	}
}

func Test_basic_infos_ScrollPgDown(t *testing.T) {
	m := makeModel()
	m.infosLinesAmount = 100

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 10, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()

	basic_infos_ScrollPgDown(m)
	if m.infosScrollOffset != 5 {
		t.Errorf("Expected 5, got %d", m.infosScrollOffset)
	}
}

func Test_min(t *testing.T) {
	if MinTestHelper(2, 3) != 2 {
		t.Error("min(2,3) should be 2")
	}
	if MinTestHelper(-1, 3) != 0 {
		t.Error("min(-1,3) should be 0")
	}
	if MinTestHelper(2, -3) != 0 {
		t.Error("min(2,-3) should be 0")
	}
	if MinTestHelper(5, 2) != 2 {
		t.Error("min(5,2) should be 2")
	}
}

func Test_basic_infos_Render_Slow(t *testing.T) {
	m := makeModel()
	m.response.ExecutionTime = 1000
	m.config.SlowResponseThreshold = 10

	out := basic_infos_Render(m)
	if !strings.Contains(out, "slow") {
		t.Error("Should indicate slow")
	}
}

func Test_basic_infos_Render_Warn(t *testing.T) {
	m := makeModel()
	m.response.ExecutionTime = 8
	m.config.SlowResponseThreshold = 10

	out := basic_infos_Render(m)
	if !strings.Contains(out, "slow") {
		t.Error("Should indicate slow (warn)")
	}
}

func Test_basic_infos_Render_Fast(t *testing.T) {
	m := makeModel()
	m.response.ExecutionTime = 1
	m.config.SlowResponseThreshold = 10

	out := basic_infos_Render(m)
	if !strings.Contains(out, "fast") {
		t.Error("Should indicate fast")
	}
}

func Test_basic_infos_Render_NoBody(t *testing.T) {
	m := makeModel()
	m.response.Body = []request_module.RequestBody{}

	out := basic_infos_Render(m)
	if strings.Contains(out, "Request Body:") {
		t.Error("Should not show body if empty")
	}
}

func Test_basic_infos_Render_Paged_NoScroll(t *testing.T) {
	m := makeModel()
	m.response.Result = strings.Repeat("line\n", 2)

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 50, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()

	out := basic_infos_Render_Paged(m)
	if strings.Contains(out, "scroll") {
		t.Error("Should not show scroll instruction if not needed")
	}
}

func Test_basic_infos_ScrollDown_NoAdvance(t *testing.T) {
	m := makeModel()
	m.infosLinesAmount = 5

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 50, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()

	m.infosScrollOffset = 0

	basic_infos_ScrollDown(m)
	if m.infosScrollOffset != 0 {
		t.Error("Should not advance scroll if there are not enough lines")
	}
}

func Test_basic_infos_ScrollPgDown_NoAdvance(t *testing.T) {
	m := makeModel()
	m.infosLinesAmount = 5

	oldGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 50, nil }
	defer func() { term_size.GetHeightFunc = oldGetHeightFunc }()

	m.infosScrollOffset = 0

	basic_infos_ScrollPgDown(m)
	if m.infosScrollOffset != 0 {
		t.Error("Should not advance PgDown if there are not enough lines")
	}
}

func Test_basic_infos_ScrollUp_NoRetreat(t *testing.T) {
	m := makeModel()
	m.infosScrollOffset = 0

	basic_infos_ScrollUp(m)

	if m.infosScrollOffset != 0 {
		t.Error("Should not retreat scroll if already at the top")
	}
}
