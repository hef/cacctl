package client

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net"
	"strings"
)

type parsedIpv6Popup struct {
	Ipv6Address net.IPNet
	Ipv6Gateway net.IP
}

func parseIpv6Popup(reader io.Reader) (*parsedIpv6Popup, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	/* doc.Find("body").Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "#text" {
			fmt.Printf(">>> (%d) >>> %s\n", i, s.Text())
		}
	})*/

	nodes := doc.Find("body").Contents().Nodes
	if len(nodes) <= 11 {
		return nil, errors.New("could not parse ipv6 popup")
	}

	ipv6AddrString := strings.TrimSpace(nodes[7].Data)

	ip, mask, err := net.ParseCIDR(ipv6AddrString)
	if err != nil {
		return nil, err
	}

	ipv6address := net.IPNet{
		IP:   ip,
		Mask: mask.Mask,
	}

	ipv6GatewayString := strings.TrimSpace(nodes[11].Data)

	gateway := net.ParseIP(ipv6GatewayString)

	form := parsedIpv6Popup{
		Ipv6Address: ipv6address,
		Ipv6Gateway: gateway,
	}
	return &form, nil
}
