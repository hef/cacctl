package client

import (
	"context"
	"errors"
	"net/http"
)

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
