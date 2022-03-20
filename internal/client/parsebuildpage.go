package client

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
)

type parsedBuildForm struct {
	infra string
	token string
}

func parseBuildPage(reader io.Reader) (*parsedBuildForm, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	form := &parsedBuildForm{}

	var ok bool
	form.infra, ok = doc.Find("button[name=infra]:contains('Build to CloudPRO v4')").First().Attr("value")
	if !ok {
		return nil, errors.New("failed to parse page")
	}
	form.token, ok = doc.Find("input[name=token]").First().Attr("value")
	if !ok {
		return nil, errors.New("failed to parse page")
	}
	return form, nil
}

type parsedBuildFormV4 struct {
	cpu       []int32
	ramMb     []int32
	storageGb []int32
	os        map[string]int32

	totalBuilds       int32
	totalCpu          int32
	totalRam          int32
	totalSsd          int32
	percentUsedBuilds int32
	percentUsedCpu    int32
	percentUsedRam    int32
	percentUsedSSd    int32

	infra       string
	token       string
	buildSubmit string
}

func parseBuildPageV4(reader io.Reader) (*parsedBuildFormV4, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	form := &parsedBuildFormV4{}

	doc.Find("form[id=build-confirm] select[name=cpu] option").Each(func(_ int, selection *goquery.Selection) {
		valueString, ok := selection.Attr("value")
		if ok {
			value, err := strconv.ParseInt(valueString, 10, 32)
			if err != nil {
				return
			}
			form.cpu = append(form.cpu, int32(value))
		}
	})

	doc.Find("form[id=build-confirm] select[name=ram] option").Each(func(_ int, selection *goquery.Selection) {
		valueString, ok := selection.Attr("value")
		if ok {
			value, err := strconv.ParseInt(valueString, 10, 32)
			if err != nil {
				return
			}
			form.ramMb = append(form.ramMb, int32(value))
		}
	})

	doc.Find("form[id=build-confirm] select[name=storage] option").Each(func(_ int, selection *goquery.Selection) {
		valueString, ok := selection.Attr("value")
		if ok {
			value, err := strconv.ParseInt(valueString, 10, 32)
			if err != nil {
				return
			}
			form.storageGb = append(form.storageGb, int32(value))
		}
	})

	form.os = make(map[string]int32)
	doc.Find("form[id=build-confirm] select[name=os] option").Each(func(_ int, selection *goquery.Selection) {
		valueString, ok := selection.Attr("value")
		if ok {
			value, err := strconv.ParseInt(valueString, 10, 32)
			if err != nil {
				return
			}
			form.os[selection.Text()] = int32(value)
		}
	})

	form.totalBuilds, form.percentUsedBuilds = parseValueFromAvailableResourcesTable(doc, " Daily Limit")
	form.totalCpu, form.percentUsedCpu = parseValueFromAvailableResourcesTable(doc, " CPU")
	form.totalRam, form.percentUsedRam = parseValueFromAvailableResourcesTable(doc, " MB RAM")
	form.totalSsd, form.percentUsedSSd = parseValueFromAvailableResourcesTable(doc, " GB SSD")

	buildSubmit, ok := doc.Find("button[name=build-submit]").First().Attr("value")
	if ok {
		form.buildSubmit = buildSubmit
	}

	infra, ok := doc.Find("input[name=infra]").First().Attr("value")
	if ok {
		form.infra = infra
	}

	token, ok := doc.Find("input[name=token]").First().Attr("value")
	if ok {
		form.token = token
	}

	return form, nil

}

func parseValueFromAvailableResourcesTable(doc *goquery.Document, resourceType string) (total, percent int32) {
	selector := fmt.Sprintf("td:contains('%s:')", resourceType)
	resourceNode := doc.Find(selector).First()
	resourceString := resourceNode.Text()
	if len(resourceString) > len(resourceType)+1 {
		resourceString = resourceString[:len(resourceString)-len(resourceType)-1]
		resource, err := strconv.ParseInt(resourceString, 10, 32)
		if err == nil {
			total = int32(resource)
		}
		resourcePercentString := strings.TrimSpace(resourceNode.Next().Text())
		if len(resourcePercentString) > 1 {
			resourcePercentString = resourcePercentString[:len(resourcePercentString)-1]
			resourcePercent, err := strconv.ParseInt(resourcePercentString, 10, 32)
			if err == nil {
				percent = int32(resourcePercent)
			}
		}
	}
	return
}
