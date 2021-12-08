package client

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	ServerName string
	ServerId   int
	Installed  time.Time
	IpAddress  net.IP
	Netmask    net.IP
	Gateway    net.IP
	Password   string

	CurrentOs string
	Ipv4      net.IP
	Ipv6      net.IP

	Hostname string
	CpuCount int
	RamMB    int
	SsdGB    int

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
	req.Header.Set("User-Agent", "cacctl/0.0.0")

	resp, err := c.c.Do(req)
	if errors.Is(err, needsLoginErr) {
		c.login(ctx)
		req, err = http.NewRequestWithContext(ctx, "GET", "https://panel.cloudatcost.com", nil)
		req.Header.Set("User-Agent", "cacctl/0.0.0")
		if err != nil {
			return nil, err
		}
		resp, err = c.c.Do(req)
	}

	if err != nil {
		return nil, err
	}

	debugPrintResp(resp, nil)

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
		currentOs := strings.TrimSpace(selection.Find("td:contains('Current Os:')").Find("td").Next().First().Text())
		server.CurrentOs = currentOs

		//serverIdString := strings.TrimSpace(selection.Find("tr:contains('Server ID:')").Find("td").Text())
		//serverId, _ := strconv.ParseInt(serverIdString, 10, 32)
		//server.ServerId = int(serverId)

		x := selection.Find("td:contains('IPv4:')").First().Text()
		log.Printf("%s", x)

		log.Println(selection.Html())

		servers = append(servers, server)

	})

	return servers, nil
}
