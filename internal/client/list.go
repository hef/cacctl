package client

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
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

	return &ListResponse{}, nil
}

func parseServersFromBody(reader io.Reader) []Server {
	return nil
}
