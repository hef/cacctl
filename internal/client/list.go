package client

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	ServerName string
	ServerId   int64
	Installed  time.Time
	IpAddress  net.IP
	Netmask    net.IP
	Gateway    net.IP
	Password   string

	CurrentOs string
	Ipv4      net.IP
	Ipv6      net.IP

	Hostname string
	CpuCount int32
	RamMB    int32
	SsdGB    int32

	Package string
}

type ListResponse struct {
	Servers []Server
}

func (c *Client) List(ctx context.Context) (*ListResponse, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", "https://panel.cloudatcost.com", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.c.Do(req)
	if errors.Is(err, needsLoginErr) {
		c.login(ctx)
		req, err = http.NewRequestWithContext(ctx, "GET", "https://panel.cloudatcost.com", nil)
		req.Header.Set("User-Agent", c.userAgent)
		if err != nil {
			return nil, err
		}
		resp, err = c.c.Do(req)
	}

	if err != nil {
		return nil, err
	}

	//debugPrintResp(resp, nil)

	servers, err := parseServersFromBody(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Servers: servers,
	}, nil
}

func parseServersFromBody(reader io.Reader) ([]Server, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var servers []Server
	doc.Find(".panel.panel-default").Each(func(i int, selection *goquery.Selection) {
		server := Server{}
		serverName := strings.TrimSpace(selection.Find("td").First().Text())
		if serverName == "" {
			return
		}
		server.ServerName = serverName
		server.CurrentOs = strings.TrimSpace(selection.Find("td td:contains('Current Os:')").Next().Text())
		server.Ipv4 = net.ParseIP(strings.TrimSpace(selection.Find("td td:contains('IPv4:')").Next().Text()))
		server.Ipv6 = net.ParseIP(strings.TrimSpace(selection.Find("td td:contains('IPv6:')").Next().Text()))

		cpuCountString := selection.Find("td td:contains(' CPU:')").First().Text()
		cpuCountString = cpuCountString[:len(cpuCountString)-len(" CPU:")]
		cpuCount, err := strconv.ParseInt(cpuCountString, 10, 32)
		if err != nil {
			return
		}
		server.CpuCount = int32(cpuCount)

		ramMBString := selection.Find("td td:contains('MB Ram:')").First().Text()
		ramMBString = ramMBString[:len(ramMBString)-len("MB Ram:")]
		ramMB, err := strconv.ParseInt(ramMBString, 10, 32)
		if err != nil {
			return
		}
		server.RamMB = int32(ramMB)

		ssdGBString := selection.Find("td td:contains('GB SSD:')").First().Text()
		ssdGBString = ssdGBString[:len(ssdGBString)-len("GB SSD:")]
		ssdGB, err := strconv.ParseInt(ssdGBString, 10, 32)
		if err != nil {
			return
		}
		server.SsdGB = int32(ssdGB)

		server.Package = strings.TrimSpace(selection.Find("[id^=Body_].panel-collapse.in div").Last().Text())

		infoBox, ok := selection.Find("[id^=Info_]").First().Attr("data-content")
		if ok {
			info, err := goquery.NewDocumentFromReader(strings.NewReader(infoBox))
			if err != nil {
				return
			}
			ServerIdString := strings.TrimSpace(info.Find("td:contains('Server ID:')").Next().Text())
			server.ServerId, err = strconv.ParseInt(ServerIdString, 10, 64)
			if err != nil {
				return
			}

			InstallDateString := strings.TrimSpace(info.Find("td:contains('Installed:')").Next().Text())
			server.Installed, err = time.Parse("01/02/2006", InstallDateString)
			if err != nil {
				return
			}
			server.IpAddress = net.ParseIP(strings.TrimSpace(info.Find("td:contains('IP Address:')").Next().Text()))
			server.Netmask = net.ParseIP(strings.TrimSpace(info.Find("td:contains('Netmask:')").Next().Text()))
			server.Gateway = net.ParseIP(strings.TrimSpace(info.Find("td:contains('Gateway:')").Next().Text()))
			server.Password = strings.TrimSpace(info.Find("td:contains('Password:')").Next().Text())
		}
		servers = append(servers, server)

	})

	return servers, nil
}
