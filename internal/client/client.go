package client

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	username string
	password string
	c        *http.Client
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

	if c.c == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		c.c = &http.Client{
			Jar: jar,
		}
		c.c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if req.URL.String() == "https://panel.cloudatcost.com/login.php" {
				return needsLoginErr
			}
			return nil
		}
	}

	return c, nil
}
