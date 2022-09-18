package client

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestParseIpv6Popup(t *testing.T) {
	f, err := os.Open("testdata/ipv6_popup.html")
	if err != nil {
		panic(err)
	}

	form, err := parseIpv6Popup(f)
	if err != nil {
		t.Errorf("error parsing build page: %s", err)
	}

	expectedAddress := "2607:8880::1830:9d/120"

	if expectedAddress != form.Ipv6Address.String() {
		t.Errorf("expected ipv6 address %s, got %s", expectedAddress, form.Ipv6Address.String())
	}

	expectedGateway := "2607:8880::1830:1"

	if expectedGateway != form.Ipv6Gateway.String() {
		t.Errorf("expected ipv6 gateway %s, got %s", expectedGateway, form.Ipv6Gateway.String())
	}

}

func FuzzParseIpv6Popup(f *testing.F) {
	c, err := os.Open("testdata/ipv6_popup.html")
	if err != nil {
		panic(err)
	}
	sample, err := io.ReadAll(c)
	if err != nil {
		panic(err)
	}
	f.Add(sample)
	f.Fuzz(func(t *testing.T, data []byte) {
		parseIpv6Popup(bytes.NewReader(data))
	})
}
