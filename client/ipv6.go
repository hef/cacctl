package client

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

type GetIpv6Response struct {
	Ipv6Address net.IPNet
	Ipv6Gateway net.IP
}

func (c *Client) GetIpv6(ctx context.Context, sid int64) (*GetIpv6Response, error) {

	q := url.Values{}
	q.Set("sid", strconv.FormatInt(sid, 10))
	u := url.URL{
		Scheme:   "https",
		Host:     "panel.cloudatcost.com",
		Path:     "/panel/_config/pop/ipv6.php",
		RawQuery: q.Encode(),
	}

	resp, err := c.sendRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	parsedIpv6, err := parseIpv6Popup(resp.Body)
	if err != nil {
		return nil, err
	}

	return &GetIpv6Response{
		Ipv6Address: parsedIpv6.Ipv6Address,
		Ipv6Gateway: parsedIpv6.Ipv6Gateway,
	}, nil
}
