package client

import (
	"context"
	"net/http"
	"net/url"
)

type LoginResponse struct {
}

func (c *Client) login(ctx context.Context) (*LoginResponse, error) {
	return c.loginWithUsernameAndPassword(ctx, c.username, c.password)
}

func (c *Client) loginWithUsernameAndPassword(ctx context.Context, username, password string) (*LoginResponse, error) {

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	form.Add("failedpage", "login-failed-2.php")
	form.Add("submit", "Login")

	req, err := c.newRequest(ctx, http.MethodPost, "https://panel.cloudatcost.com/manage-check2.php", &form)
	_, err = c.httpClient.Do(req) // todo: check response for failed password
	if err != nil {
		return nil, err
	}

	return &LoginResponse{}, nil

}
