package client

import "github.com/PuerkitoBio/goquery"

func parseIpv6Popup(reader io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	var ipv6 string
	doc.Find("input[name=ipv6]").First().Each(func(i int, s *goquery.Selection) {
		ipv6, _ = s.Attr("value")
	})
	return ipv6, nil
}
