package client

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type LoginResponse struct {
}

func (c *Client) login(ctx context.Context) (*LoginResponse, error) {
	return c.loginWithUsernameAndPassword(ctx, c.username, c.password)
}

func (c *Client) loginWithUsernameAndPassword(ctx context.Context, username, password string) (*LoginResponse, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://panel.cloudatcost.com/login.php", nil)
	req.Header.Set("User-Agent", "cacctl/0.0.0")
	if err != nil {
		return nil, err
	}
	c.c.Do(req)

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	form.Add("failedpage", "login-failed-2.php")
	form.Add("submit", "Login")

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, "https://panel.cloudatcost.com/manage-check2.php", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "cacctl/0.0.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = c.c.Do(req)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{}, nil

}
