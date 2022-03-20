package client

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strconv"
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

	VmName     string
	CustomerId int64

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

// <option value='1'>All</option>
// <option value='2'>Powered On</option>
// <option value='3'>Powered Off</option>
// <option value='4'>Installing</option>
// <option value='5'>Pending</option>
// <option value='6'>Installed</option>
type ListFilter int

const (
	All        ListFilter = 1
	PoweredOn             = 2
	PoweredOff            = 3
	Installing            = 4
	Pending               = 5
	Installed             = 6
)

func (f ListFilter) asInt() int {
	return int(f)
}

func ListFilterFromString(filter string) ListFilter {
	switch filter {
	case "All":
		return All
	case "PoweredOn":
		return PoweredOn
	case "PoweredOff":
		return PoweredOff
	case "Installing":
		return Installing
	case "Pending":
		return Pending
	case "Installed":
		return Installed
	}
	return 0 // this is probably a bad thing to do
}

func (c *Client) ListWithFilter(ctx context.Context, search string, limit int, filter ListFilter) (*ListResponse, error) {

	form := url.Values{}
	form.Add("search", search)
	form.Add("limit", strconv.Itoa(limit))
	form.Add("filter", strconv.Itoa(filter.asInt()))

	resp, err := c.sendRequest(ctx, http.MethodPost, "https://panel.cloudatcost.com", &form)
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
