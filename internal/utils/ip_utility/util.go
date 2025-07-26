package ip_utility

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	ip_cache_module "github.com/diogopereiradev/httpzen/internal/cache"
	"github.com/go-resty/resty/v2"
)

type LookupIpInfo struct {
	Type      string  `json:"type"`
	Ip        string  `json:"ip"`
	Decimal   string  `json:"decimal,omitempty"`
	Hostname  string  `json:"hostname,omitempty"`
	ASN       string  `json:"asn,omitempty"`
	ISP       string  `json:"isp,omitempty"`
	State     string  `json:"state,omitempty"`
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

func FetchIpInfo(ipType string, ip string) (LookupIpInfo, error) {
	key := strings.ReplaceAll(ip, ".", "_")
	if cached, ok := ip_cache_module.GetIpInfoFromCache(key); ok {
		var info LookupIpInfo
		b, _ := json.Marshal(cached)
		_ = json.Unmarshal(b, &info)
		return info, nil
	}

	apiUrl := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	client := resty.New()
	res, _ := client.R().Get(apiUrl)

	type ipinfoResponse struct {
		Hostname string `json:"hostname"`
		Postal   string `json:"postal"`
		Org      string `json:"org"`
		City     string `json:"city"`
		Region   string `json:"region"`
		Country  string `json:"country"`
		Loc      string `json:"loc"`
		ASN      string `json:"asn"`
		ISP      string `json:"isp"`
	}

	var resp ipinfoResponse
	json.Unmarshal(res.Body(), &resp)

	latitude := 0.0
	longitude := 0.0

	if resp.Loc != "" {
		parts := strings.Split(resp.Loc, ",")
		if len(parts) == 2 {
			fmt.Sscanf(parts[0], "%f", &latitude)
			fmt.Sscanf(parts[1], "%f", &longitude)
		}
	}

	info := LookupIpInfo{
		Type:      ipType,
		Ip:        ip,
		Decimal:   resp.Postal,
		Hostname:  resp.Hostname,
		ASN:       resp.Org,
		ISP:       resp.Org,
		City:      resp.City,
		Country:   resp.Country,
		State:     resp.Region,
		Latitude:  latitude,
		Longitude: longitude,
	}

	ip_cache_module.SetIpInfoToCache(key, map[string]any{
		"Type":      info.Type,
		"Ip":        info.Ip,
		"Decimal":   info.Decimal,
		"Hostname":  info.Hostname,
		"ASN":       info.ASN,
		"ISP":       info.ISP,
		"City":      info.City,
		"Country":   info.Country,
		"State":     info.State,
		"Latitude":  info.Latitude,
		"Longitude": info.Longitude,
	})
	return info, nil
}

func LookupDomainIps(res *resty.Response) []LookupIpInfo {
	host := res.Request.RawRequest.URL.Host
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return []LookupIpInfo{}
	}

	var ipList []LookupIpInfo
	for _, ip := range ips {
		ipType := ""
		if ip.To4() != nil {
			ipType = "IPv4"
		} else {
			ipType = "IPv6"
		}

		ipInfo, err := FetchIpInfo(ipType, ip.String())
		if err != nil {
			continue
		}

		ipList = append(ipList, ipInfo)
	}
	return ipList
}