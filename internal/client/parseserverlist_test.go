package client

import (
	"net"
	"os"
	"testing"
	"time"
)

func mustParseDate(value string) time.Time {
	date, err := time.Parse("01/02/2006", value)
	if err != nil {
		panic(err)
	}
	return date
}

func TestParseServerList(t *testing.T) {

	var ListTests = []struct {
		Name            string
		Fixture         string
		ExpectedServers []Server
	}{
		{
			Name:    "Basic",
			Fixture: "testdata/login_success.html",
			ExpectedServers: []Server{{
				ServerName: "c999963378-cloudpro-799792359",
				ServerId:   255330174,
				Status:     "Powered On",
				Installed:  mustParseDate("12/06/2021"),
				IpAddress:  net.IPv4(142, 47, 88, 68),
				Netmask:    net.IPv4(255, 255, 255, 0),
				Gateway:    net.IPv4(142, 47, 88, 1),
				Password:   "3aXz7qSd9S",
				VmName:     "c999963378-CloudPRO-519046902-629183859",
				CustomerId: 999963378,
				CurrentOs:  "Ubuntu 18.04 LTS 64bit",
				Ipv4:       net.IPv4(142, 47, 88, 68),
				Ipv6:       net.ParseIP("2607:8880::A000:1044"),
				Hostname:   "",
				CpuCount:   2,
				RamMB:      4096,
				SsdGB:      40,
				Package:    "CloudPRO v4",
			}},
		},
		{
			Name:    "InstallingPair",
			Fixture: "testdata/2_installing.html",
			ExpectedServers: []Server{{
				ServerName: "", //todo: "c999963378-cloudpro-313939255",
				ServerId:   255330530,
				Status:     "Installing",
				Installed:  mustParseDate("12/08/2021"),
				IpAddress:  net.IPv4(142, 47, 88, 224),
				Netmask:    net.IPv4(255, 255, 255, 0),
				Gateway:    net.IPv4(142, 47, 88, 1),
				Password:   "n22rW2BG4d",
				CustomerId: 0, // Not reliable source.  Could extract it from server name maybe.
				CurrentOs:  "Ubuntu 18.04 LTS 64bit",
				Ipv4:       net.IPv4(142, 47, 88, 224),
				Ipv6:       net.ParseIP("2607:8880::a000:10e0"),
				Hostname:   "",
				CpuCount:   1,
				RamMB:      2048,
				SsdGB:      20,
				Package:    "CloudPRO v4",
			}, {
				ServerName: "", // todo: "c999963378-cloudpro-769040999",
				ServerId:   255330529,
				Status:     "Installing",
				Installed:  mustParseDate("12/08/2021"),
				IpAddress:  net.IPv4(142, 47, 88, 223),
				Netmask:    net.IPv4(255, 255, 255, 0),
				Gateway:    net.IPv4(142, 47, 88, 1),
				CustomerId: 0, // Not reliable source.  Could extract it from server name maybe.
				Password:   "3YAqeQPz6t",
				CurrentOs:  "Ubuntu 18.04 LTS 64bit",
				Ipv4:       net.IPv4(142, 47, 88, 223),
				Ipv6:       net.ParseIP("2607:8880::a000:10df"),
				Hostname:   "",
				CpuCount:   1,
				RamMB:      2048,
				SsdGB:      20,
				Package:    "CloudPRO v4",
			}},
		},
	}

	for _, tt := range ListTests {

		t.Run(tt.Name, func(t *testing.T) {

			f, err := os.Open(tt.Fixture)
			if err != nil {
				panic(err)
			}
			servers, err := parseServersFromBody(f)
			if err != nil {
				t.Errorf("unexpected error parsing body: %s", err)
			}

			if len(servers) != len(tt.ExpectedServers) {
				t.Errorf("expected Server Count %d, got %d", len(tt.ExpectedServers), len(servers))
			}

			for i, server := range servers {
				expectedServer := tt.ExpectedServers[i]
				if server.ServerName != expectedServer.ServerName {
					t.Errorf("expected ServerName %s, got %s", expectedServer.ServerName, server.ServerName)
				}
				if server.ServerId != expectedServer.ServerId {
					t.Errorf("expected ServerId %d, got %d", expectedServer.ServerId, server.ServerId)
				}
				if server.Status != expectedServer.Status {
					t.Errorf("expected Status %s, got %s", expectedServer.Status, server.Status)
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
				if expectedServer.VmName != server.VmName {
					t.Errorf("expected VM Name %s, got %s", expectedServer.VmName, server.VmName)
				}

				if expectedServer.CustomerId != server.CustomerId {
					t.Errorf("expected CustomerID %d, got %d", expectedServer.CustomerId, server.CustomerId)
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
		})
	}
}

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
