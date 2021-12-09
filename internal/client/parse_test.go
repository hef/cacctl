package client

import (
	"net"
	"os"
	"testing"
	"time"
)

func TestParseLogin(t *testing.T) {
	f, err := os.Open("testdata/login_success.html")
	if err != nil {
		panic(err)
	}

	installDate, _ := time.Parse("01/02/2006", "12/06/2021")

	expectedServer := Server{
		ServerName: "c999963378-cloudpro-799792359",
		ServerId:   255330174,
		Installed:  installDate,
		IpAddress:  net.IPv4(142, 47, 88, 68),
		Netmask:    net.IPv4(255, 255, 255, 0),
		Gateway:    net.IPv4(142, 47, 88, 1),
		Password:   "3aXz7qSd9S",
		CurrentOs:  "Ubuntu 18.04 LTS 64bit",
		Ipv4:       net.IPv4(142, 47, 88, 68),
		Ipv6:       net.ParseIP("2607:8880::A000:1044"),
		Hostname:   "",
		CpuCount:   2,
		RamMB:      4096,
		SsdGB:      40,
		Package:    "CloudPRO v4",
	}

	servers, err := parseServersFromBody(f)
	if err != nil {
		t.Errorf("expected no error from ParseServersFrom body, got %s", err)
	}

	if len(servers) != 1 {
		t.Errorf("got unexpected number of servers back, expected 1, got %d", len(servers))
	}
	if len(servers) > 0 {
		server := servers[0]
		if server.ServerName != expectedServer.ServerName {
			t.Errorf("expected ServerName %s, got %s", expectedServer.ServerName, server.ServerName)
		}
		if server.ServerId != expectedServer.ServerId {
			t.Errorf("expected ServerId %d, got %d", expectedServer.ServerId, server.ServerId)
		}
		if server.Installed != expectedServer.Installed {
			t.Errorf("expected Install Date %s, got %s", expectedServer.Installed, server.Installed)
		}
		if !expectedServer.IpAddress.Equal(server.IpAddress) {
			t.Errorf("expected IP Address %s, got %s", expectedServer.IpAddress, server.IpAddress)
		}
		if !expectedServer.Netmask.Equal(server.Netmask) {
			t.Errorf("expected Netmask %s, got %s", expectedServer.Netmask, server.Netmask)
		}
		if !expectedServer.Gateway.Equal(server.Gateway) {
			t.Errorf("expected Gateway %s, got %s", expectedServer.Gateway, server.Gateway)
		}
		if expectedServer.Password != server.Password {
			t.Errorf("expected Password %s, got %s", expectedServer.Password, server.Password)
		}
		if expectedServer.CurrentOs != server.CurrentOs {
			t.Errorf("expected Current OS: %s, got %s", expectedServer.CurrentOs, server.CurrentOs)
		}
		if !expectedServer.Ipv4.Equal(server.Ipv4) {
			t.Errorf("expected Ipv4: %s, got %s", expectedServer.Ipv4, server.Ipv4)
		}
		if !expectedServer.Ipv6.Equal(server.Ipv6) {
			t.Errorf("expected Ipv6: %s, got %s", expectedServer.Ipv6, server.Ipv6)
		}
		if expectedServer.Hostname != server.Hostname {
			t.Errorf("expected Hostname: %s, got %s", expectedServer.Hostname, server.Hostname)
		}
		if expectedServer.CpuCount != server.CpuCount {
			t.Errorf("expected cpu count %d, got %d", expectedServer.CpuCount, server.CpuCount)
		}
		if expectedServer.RamMB != server.RamMB {
			t.Errorf("expected Ram (MB) %d, got %d", expectedServer.RamMB, server.RamMB)
		}
		if expectedServer.SsdGB != server.SsdGB {
			t.Errorf("expected Ssd (GB) %d, got %d", expectedServer.SsdGB, server.SsdGB)
		}
		if expectedServer.Package != server.Package {
			t.Errorf("expected Package %s, got %s", expectedServer.Package, server.Package)
		}
	}
}
