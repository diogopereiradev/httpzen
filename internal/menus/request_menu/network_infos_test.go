package request_menu

import (
	"testing"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/diogopereiradev/httpzen/internal/utils/term_size"
	"github.com/stretchr/testify/assert"
)

func newModelWithIpInfos(ipInfos []request_module.IpInfo, scrollOffset int) *Model {
	return &Model{
		response:            &request_module.RequestResponse{IpInfos: ipInfos},
		networkScrollOffset: scrollOffset,
	}
}

func Test_network_infos_Render(t *testing.T) {
	infos := []request_module.IpInfo{{
		Type: "IPv4", Ip: "1.2.3.4", Country: "BR", City: "SP", Decimal: "1234", Hostname: "host",
		State: "SP", ASN: "AS123", ISP: "ISP", Latitude: 1.23, Longitude: 4.56,
	}}
	m := newModelWithIpInfos(infos, 0)
	out := network_infos_Render(m)
	assert.Contains(t, out, "IPv4")
	assert.Contains(t, out, "1.2.3.4")
	assert.Contains(t, out, "Country: BR")
	assert.Contains(t, out, "City: SP")
	assert.Contains(t, out, "Decimal: 1234")
	assert.Contains(t, out, "Hostname: host")
	assert.Contains(t, out, "Region/State: SP")
	assert.Contains(t, out, "ASN: AS123")
	assert.Contains(t, out, "ISP: ISP")
	assert.Contains(t, out, "Coordinates: 1.23, 4.56")
}

func Test_network_infos_Render_empty(t *testing.T) {
	m := newModelWithIpInfos(nil, 0)
	out := network_infos_Render(m)
	assert.Equal(t, "", out)
}

func Test_network_infos_Render_Paged(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	infos := []request_module.IpInfo{
		{Type: "IPv4", Ip: "1.2.3.4"},
		{Type: "IPv6", Ip: "::1"},
	}
	m := newModelWithIpInfos(infos, 0)
	out := network_infos_Render_Paged(m)
	assert.Contains(t, out, "IPv4")
	assert.NotContains(t, out, "IPv6")
	assert.Contains(t, out, "[1-2/")

	m.networkScrollOffset = 3
	out2 := network_infos_Render_Paged(m)
	assert.Contains(t, out2, "IPv6")
}

func Test_network_infos_ScrollUp(t *testing.T) {
	m := newModelWithIpInfos(nil, 2)
	network_infos_ScrollUp(m)
	assert.Equal(t, 1, m.networkScrollOffset)
	network_infos_ScrollUp(m)
	assert.Equal(t, 0, m.networkScrollOffset)
	network_infos_ScrollUp(m)
	assert.Equal(t, 0, m.networkScrollOffset)
}

func Test_network_infos_ScrollDown(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithIpInfos([]request_module.IpInfo{{}, {}}, 0)
	m.networkLinesAmount = 3
	network_infos_ScrollDown(m)
	assert.Equal(t, 1, m.networkScrollOffset)
	network_infos_ScrollDown(m)
	assert.Equal(t, 1, m.networkScrollOffset)
}

func Test_network_infos_ScrollDown_borders(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 20, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithIpInfos(nil, 0)
	m.networkLinesAmount = 0
	network_infos_ScrollDown(m)
	assert.Equal(t, 0, m.networkScrollOffset)
	m.networkLinesAmount = 2
	network_infos_ScrollDown(m)
	assert.Equal(t, 0, m.networkScrollOffset)
}

func Test_network_infos_ScrollPgUp(t *testing.T) {
	m := newModelWithIpInfos(nil, 6)
	network_infos_ScrollPgUp(m)
	assert.Equal(t, 1, m.networkScrollOffset)
	network_infos_ScrollPgUp(m)
	assert.Equal(t, 0, m.networkScrollOffset)
}

func Test_network_infos_ScrollPgDown(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 18, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithIpInfos([]request_module.IpInfo{{}, {}, {}}, 0)
	m.networkLinesAmount = 3
	network_infos_ScrollPgDown(m)
	assert.Equal(t, 5, m.networkScrollOffset)
}

func Test_network_infos_ScrollPgDown_borders(t *testing.T) {
	origGetHeightFunc := term_size.GetHeightFunc
	term_size.GetHeightFunc = func() (int, error) { return 20, nil }
	defer func() { term_size.GetHeightFunc = origGetHeightFunc }()

	m := newModelWithIpInfos(nil, 0)
	m.networkLinesAmount = 0
	network_infos_ScrollPgDown(m)
	assert.Equal(t, 0, m.networkScrollOffset)
	m.networkLinesAmount = 2
	network_infos_ScrollPgDown(m)
	assert.Equal(t, 0, m.networkScrollOffset)
}
