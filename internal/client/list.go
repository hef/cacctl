package client

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	ServerName string
	ServerId   int64
	Status     string
	Installed  time.Time
	IpAddress  net.IP
	Netmask    net.IP
	Gateway    net.IP
	Password   string

	VmName string

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

	resp, err := c.sendRequest(ctx, http.MethodGet, "https://panel.cloudatcost.com", nil)
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
