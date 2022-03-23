package client

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// PowerCycle(2, "c999963378-CloudPRO-519046902-629183859", "255330174", "c999963378-cloudpro-799792359")
var powerCycleCallRegex = regexp.MustCompile(`PowerCycle\(\d, "(?P<vmname>[\w\-]+)", "\d+", "[\w\-]+"\)`)

// DELETECPRO2(255330174, "c999963378-cloudpro-799792359", "999963378", "c999963378-CloudPRO-519046902-629183859", "v4")
var deleteCPro2CallRegex = regexp.MustCompile(`DELETECPRO2\((?P<sid>\d+), "(?P<servername>[\w\-]+)", "(?P<cid>\d+)", "(?P<vmname>[\w\-]+)", "v4"\)`)

func parseServersFromBody(reader io.Reader) ([]Server, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var servers []Server
	doc.Find(".panel.panel-default").Each(func(i int, selection *goquery.Selection) {
		server := Server{}
		TitleNode := selection.Find("td[id^=PanelTitle_]").First().Text()
		if TitleNode == "" {
			return
		} else if TitleNode == "  Install Pending  \n\t\t\t\t\t  Status " {
			server.Status = "Install Pending"
		} else if TitleNode == "  Installing " {
			server.Status = "Installing"
		} else if TitleNode == "  Install Failed  \n\t\t\t\t\t  Details" {
			server.Status = "Install Failed"
		} else {
			server.ServerName = strings.TrimSpace(selection.Find("td[id^=PanelTitle_]").First().Text())
		}

		icon, ok := selection.Find("td[id^=PanelTitle_] i").First().Attr("class")
		if ok {
			if icon == "fa fa-cloud-upload" {
				server.Status = "Powered On"
			}
		}

		server.CurrentOs = strings.TrimSpace(selection.Find("td td:contains('Current Os:')").Next().Text())
		server.Ipv4 = net.ParseIP(strings.TrimSpace(selection.Find("td td:contains('IPv4:')").Next().Text()))
		server.Ipv6 = net.ParseIP(strings.TrimSpace(selection.Find("td td:contains('IPv6:')").Next().Text()))

		cpuCountString := selection.Find("td td:contains(' CPU:')").First().Text()
		cpuCountString = cpuCountString[:len(cpuCountString)-len(" CPU:")]
		cpuCount, err := strconv.ParseInt(cpuCountString, 10, 32)
		if err != nil {
			return
		}
		server.CpuCount = int32(cpuCount)

		ramMBString := selection.Find("td td:contains('MB Ram:')").First().Text()
		ramMBString = ramMBString[:len(ramMBString)-len("MB Ram:")]
		ramMB, err := strconv.ParseInt(ramMBString, 10, 32)
		if err != nil {
			return
		}
		server.RamMB = int32(ramMB)

		ssdGBString := selection.Find("td td:contains('GB SSD:')").First().Text()
		ssdGBString = ssdGBString[:len(ssdGBString)-len("GB SSD:")]
		ssdGB, err := strconv.ParseInt(ssdGBString, 10, 32)
		if err != nil {
			return
		}
		server.SsdGB = int32(ssdGB)

		server.Package = strings.TrimSpace(selection.Find("[id^=Body_].panel-collapse div").Last().Text())

		infoBox, ok := selection.Find("[id^=Info_]").First().Attr("data-content")
		if ok {
			info, err := goquery.NewDocumentFromReader(strings.NewReader(infoBox))
			if err != nil {
				return
			}
			ServerIdString := strings.TrimSpace(info.Find("td:contains('Server ID:')").Next().Text())
			server.ServerId, err = strconv.ParseInt(ServerIdString, 10, 64)
			if err != nil {
				return
			}

			InstallDateString := strings.TrimSpace(info.Find("td:contains('Installed:')").Next().Text())
			server.Installed, err = time.Parse("01/02/2006", InstallDateString)
			if err != nil {
				return
			}
			server.IpAddress = net.ParseIP(strings.TrimSpace(info.Find("td:contains('IP Address:')").Next().Text()))
			server.Netmask = net.ParseIP(strings.TrimSpace(info.Find("td:contains('Netmask:')").Next().Text()))
			server.Gateway = net.ParseIP(strings.TrimSpace(info.Find("td:contains('Gateway:')").Next().Text()))
			server.Password = strings.TrimSpace(info.Find("td:contains('Password:')").Next().Text())

		}
		PowerCycleFunctionCall, ok := selection.Find("a[onclick*='PowerCycle']").First().Attr("onclick")
		if ok {
			m := powerCycleCallRegex.FindStringSubmatch(PowerCycleFunctionCall)
			if len(m) > 0 {
				server.VmName = m[1]
			}
		}

		deleteCPro2Call, ok := selection.Find("a[onclick*='DELETECPRO2']").First().Attr("onclick")
		if ok {
			m := deleteCPro2CallRegex.FindStringSubmatch(deleteCPro2Call)
			if len(m) > 3 {
				//server.VmName = m[1]
				server.CustomerId, _ = strconv.ParseInt(m[3], 10, 64)
			}
		}
		servers = append(servers, server)
	})

	return servers, nil
}
