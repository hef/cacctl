package client

import (
	"errors"
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
	form.infra, ok = doc.Find("button[name=infra]").First().Attr("value")
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

	{
		dailyLimitNode := doc.Find("td:contains(' Daily Limit:')").First()
		dailyLimitString := dailyLimitNode.Text()
		dailyLimitString = dailyLimitString[:len(dailyLimitString)-len(" Daily Limit:")]
		dailyLimit, err := strconv.ParseInt(dailyLimitString, 10, 32)
		if err == nil {
			form.totalBuilds = int32(dailyLimit)
		}
		DailyLimitPercentString := strings.TrimSpace(dailyLimitNode.Next().Text())
		DailyLimitPercentString = DailyLimitPercentString[:len(DailyLimitPercentString)-1]
		dailyLimitPercent, err := strconv.ParseInt(DailyLimitPercentString, 10, 32)
		form.percentUsedBuilds = int32(dailyLimitPercent)
	}

	{
		cpuNode := doc.Find("td:contains(' CPU:')").First()
		cpuString := cpuNode.Text()
		cpuString = cpuString[:len(cpuString)-len(" CPU:")]
		cpu, err := strconv.ParseInt(cpuString, 10, 32)
		if err == nil {
			form.totalCpu = int32(cpu)
		}
		cpuPercentString := strings.TrimSpace(cpuNode.Next().Text())
		cpuPercentString = cpuPercentString[:len(cpuPercentString)-1]
		cpuPercent, err := strconv.ParseInt(cpuPercentString, 10, 32)
		form.percentUsedCpu = int32(cpuPercent)
	}

	{
		ramNode := doc.Find("td:contains(' MB RAM:')").First()
		ramString := ramNode.Text()
		ramString = ramString[:len(ramString)-len(" MB RAM:")]
		ram, err := strconv.ParseInt(ramString, 10, 32)
		if err == nil {
			form.totalRam = int32(ram)
		}
		ramPercentString := strings.TrimSpace(ramNode.Next().Text())
		ramPercentString = ramPercentString[:len(ramPercentString)-1]
		ramPercent, err := strconv.ParseInt(ramPercentString, 10, 32)
		form.percentUsedRam = int32(ramPercent)
	}

	{
		ssdNode := doc.Find("td:contains(' GB SSD:')").First()
		ssdString := ssdNode.Text()
		ssdString = ssdString[:len(ssdString)-len(" GB SSD:")]
		ssd, err := strconv.ParseInt(ssdString, 10, 32)
		if err == nil {
			form.totalSsd = int32(ssd)
		}
		ssdPercentString := strings.TrimSpace(ssdNode.Next().Text())
		ssdPercentString = ssdPercentString[:len(ssdPercentString)-1]
		ssdPercent, err := strconv.ParseInt(ssdPercentString, 10, 32)
		form.percentUsedSSd = int32(ssdPercent)
	}

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
