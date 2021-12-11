package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type BuildRequest struct {
}

type BuildResponse struct {
}

func (c *Client) Build(ctx context.Context, request *BuildRequest) (*BuildResponse, error) {

	resp, err := c.sendRequest(ctx, http.MethodGet, "https://panel.cloudatcost.com/build", nil)
	if err != nil {
		return nil, err
	}

	formData, err := parseBuildPage(resp.Body)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("infra", formData.infra)
	form.Add("token", formData.token)

	resp, err = c.sendRequest(ctx, http.MethodPost, "https://panel.cloudatcost.com/build", &form)
	if err != nil {
		return nil, err
	}

	buildResponse, err := parseBuildPageV4(resp.Body)
	if err != nil {
		return nil, err
	}

	form = url.Values{}
	form.Add("build-submit", buildResponse.buildSubmit)
	form.Add("cpu", strconv.Itoa(int(buildResponse.cpu[len(buildResponse.cpu)-1])))
	form.Add("ram", strconv.Itoa(int(buildResponse.ramMb[len(buildResponse.ramMb)-1])))
	form.Add("storage", strconv.Itoa(int(buildResponse.storageGb[len(buildResponse.storageGb)-1])))
	form.Add("bs", "0")
	form.Add("ipAddress", "0")
	form.Add("ha", "0")
	form.Add("encryption", "0")
	form.Add("os", "146")
	form.Add("infra", buildResponse.infra)
	form.Add("token", buildResponse.token)

	resp, err = c.sendRequest(ctx, http.MethodPost, "https://panel.cloudatcost.com/build", &form)

	return &BuildResponse{}, nil
}
