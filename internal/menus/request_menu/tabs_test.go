package request_menu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_tab_MoveLeft(t *testing.T) {
	m := &Model{activeTab: tab_RequestInfos, resultScrollOffset: 5}
	m2 := tab_MoveLeft(m)
	assert.Equal(t, tab_Result, m2.activeTab)
	assert.Equal(t, 0, m2.resultScrollOffset)

	m = &Model{activeTab: tab_Result, resultScrollOffset: 2}
	m2 = tab_MoveLeft(m)
	assert.Equal(t, tab_ResponseHeaders, m2.activeTab)
	assert.Equal(t, 0, m2.resultScrollOffset)
}

func Test_tab_MoveRight(t *testing.T) {
	m := &Model{activeTab: tab_RequestInfos, resultScrollOffset: 3}
	m2 := tab_MoveRight(m)
	assert.Equal(t, tab_NetworkInfos, m2.activeTab)
	assert.Equal(t, 0, m2.resultScrollOffset)

	m = &Model{activeTab: tab_ResponseHeaders, resultScrollOffset: 1}
	m2 = tab_MoveRight(m)
	assert.Equal(t, tab_Result, m2.activeTab)
	assert.Equal(t, 0, m2.resultScrollOffset)
}

func Test_tab_Render(t *testing.T) {
	m := &Model{activeTab: tab_RequestHeaders}
	output := tab_Render(m)
	for i, name := range tabNames {
		assert.Contains(t, output, name, "tab name %d missing", i)
	}
}

func Test_tabBorders(t *testing.T) {
	assert.NotEqual(t, activeTabBorder, tabBorder)
	assert.Equal(t, "─", tabBorder.Top)
	assert.Equal(t, " ", activeTabBorder.Bottom)
}

func Test_tabNames(t *testing.T) {
	assert.Equal(t, 5, len(tabNames))
	assert.Equal(t, "Response", tabNames[0])
	assert.Equal(t, "Response Headers", tabNames[4])
}
