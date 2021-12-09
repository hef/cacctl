//go:build dev

package client

import "context"

type BuildRequest struct {
}

type BuildResponse struct {
}

func (c *Client) Build(ctx context.Context, request *BuildRequest) (*BuildResponse, error) {

	panic("Not Implemented")

	/*form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	form.Add("failedpage", "login-failed-2.php")
	form.Add("submit", "Login")*/

	//return &BuildResponse{}, nil
}
