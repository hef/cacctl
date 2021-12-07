package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://panel.cloudatcost.com/manage-check2.php", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "cacctl/0.0.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}


	io.Copy(os.Stdout, resp.Body)




}
