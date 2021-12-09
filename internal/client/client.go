package client

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Client struct {
	username   string
	password   string
	userAgent  string
	httpClient *http.Client
}

type err string

func (e err) Error() string {
	return string(e)
}

var (
	needsLoginErr err = "Needs Login"
)

func debugPrintResp(resp *http.Response, out io.Writer) {
	log.Printf("(%d) %s\n", resp.StatusCode, resp.Status)
	for header, values := range resp.Header {
		log.Printf("  %s: %s\n", header, values)
	}
	if out != nil {
		io.Copy(out, resp.Body)
	}
}

func New(options ...Option) (*Client, error) {
	c := &Client{}

	for _, o := range options {
		o(c)
	}

	if c.httpClient == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		c.httpClient = &http.Client{
			Jar: jar,
		}
		c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if req.URL.String() == "https://panel.cloudatcost.com/login.php" {
				return needsLoginErr
			}
			return nil
		}
	}

	if c.userAgent == "" {
		c.userAgent = "cacctl/unknown"
	}

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, method, url string, form *url.Values) (*http.Request, error) {

	var reader io.Reader
	if form != nil {
		reader = strings.NewReader(form.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	if form != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

func (c *Client) do(ctx context.Context, method, url string, form *url.Values) (*http.Response, error) {

	req, err := c.newRequest(ctx, method, url, form)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if errors.Is(err, needsLoginErr) {
		c.login(ctx)
		req, err = c.newRequest(ctx, method, url, form)
		if err != nil {
			return nil, err
		}
		resp, err = c.httpClient.Do(req)
	}
	return resp, err
}
