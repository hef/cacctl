package client

import (
	"log"
	"net/http"
)

type Client struct {
	username string
	password string
	c *http.Client
}

type err string

func (e err) Error() string {
	return string(e)
}

var (
	needsLoginErr err = "Needs Login"
)

func debugPrintResp(resp *http.Response)  {
	log.Printf("(%d) %s\n", resp.StatusCode, resp.Status)
	for header, values := range resp.Header {
		log.Printf("  %s: %s\n", header, values)
	}

}

func New(options ...Option) (*Client, error) {
	c := &Client{}

	for _, o := range options {
		o(c)
	}

	if c.c == nil {
		c.c = &http.Client{}
		c.c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if req.URL.String() == "https://panel.cloudatcost.com/login.php" {
				return needsLoginErr
			}
			return nil
		}
	}

	return c, nil
}
