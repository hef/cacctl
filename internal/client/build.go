package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type BuildRequest struct {
	Cpu     int
	Ram     int
	Storage int
	//BootScript       string
	HighAvailability bool
	Encryption       bool
	OS               string
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
	err = validateAndSet(&form, "cpu", int32(request.Cpu), buildResponse.cpu)
	if err != nil {
		return nil, err
	}
	err = validateAndSet(&form, "ram", int32(request.Ram), buildResponse.ramMb)
	if err != nil {
		return nil, err
	}
	err = validateAndSet(&form, "storage", int32(request.Storage), buildResponse.storageGb)
	if err != nil {
		return nil, err
	}
	form.Add("bs", "0")
	form.Add("ipAddress", "0")
	if request.HighAvailability {
		form.Add("ha", "1")
	} else {
		form.Add("ha", "0")
	}
	if request.Encryption {
		form.Add("encryption", "1")
	} else {
		form.Add("encryption", "0")
	}

	if request.OS == "" {
		form.Add("os", "146") // Ubuntu 18.04 LTS 64bit
	} else {
		osId, ok := buildResponse.os[request.OS]
		if !ok {
			return nil, errors.New("invalid operating system value specified") // todo: report valid values
		} else {
			form.Add("os", strconv.Itoa(int(osId)))
		}
	}

	form.Add("infra", buildResponse.infra)
	form.Add("token", buildResponse.token)

	resp, err = c.sendRequest(ctx, http.MethodPost, "https://panel.cloudatcost.com/build", &form)
	if err != nil {
		return nil, err
	}
	return &BuildResponse{}, nil
}

func validateAndSet(form *url.Values, resourceName string, requested int32, options []int32) error {
	if requested == 0 {
		if len(options) > 0 {
			form.Add(resourceName, strconv.Itoa(int(options[len(options)-1])))
			return nil
		} else {
			return fmt.Errorf("no remaining %s to allocate", resourceName)
		}
	} else {
		for _, i := range options {
			if requested == i {
				form.Add(resourceName, strconv.Itoa(int(requested)))
				return nil
			}
		}
		return fmt.Errorf("invalid %s value specified", resourceName)
	}
}
