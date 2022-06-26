package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type PowerCycle int

const (
	PowerDown PowerCycle = 0
	PowerUp              = 1
	Reboot               = 2
)

func (p PowerCycle) int() int {
	switch p {
	case PowerDown:
		return 0
	case PowerUp:
		return 1
	case Reboot:
		return 2
	default:
		panic("invalid PowerCycle value")
	}
}

func (c *Client) PowerCycle(ctx context.Context, cycle PowerCycle, vmName string, sid int64) error {

	q := url.Values{}
	q.Set("sid", strconv.FormatInt(sid, 10))
	q.Set("vmname", vmName)
	q.Set("cycle", strconv.Itoa(cycle.int()))

	u := url.URL{
		Scheme: "https",
		Host:   "panel.cloudatcost.com",
		Path:   "/panel/_config/powerCycle.php",
		//RawQuery: fmt.Sprintf("sid=%d&vmname=%s&cycle=%d", sid, vmName, cycle.int()), // order matters
		RawQuery: q.Encode(),
	}

	resp, err := c.sendRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if string(body) != "Server Successfully Rebooted" {
			return fmt.Errorf("powercycle response: %s", string(body))
		}
	} else {
		return fmt.Errorf("powercycle error: (%d) %s", resp.StatusCode, resp.Status)
	}
	return nil
}
