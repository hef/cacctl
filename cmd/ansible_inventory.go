package cmd

import (
	"context"
	"fmt"
	"github.com/hef/cacctl/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"text/template"
)

func init() {
	rootCmd.AddCommand(ansibleInventoryCmd)
	ansibleInventoryCmd.PersistentFlags().IntP("masters", "m", 3, "Number of Masters, default 3")
	viper.BindPFlag("masters", ansibleInventoryCmd.PersistentFlags().Lookup("masters"))
}

var hostsTemplate = `---
kubernetes:
  vars:
    ansible_user: root
  children:
    master:
      hosts:
        {{- range .Masters }}
        cac-{{.ServerId}}:
          ansible_host: {{.IpAddress}}
          cac_netmask: {{.Netmask}}
          cac_gateway: {{.Gateway}}
          cac_ipv6_address: {{.Ipv6Address}}
          cac_ipv6_gateway: {{.Ipv6Gateway}}
        {{- end}}
    worker:
      hosts:
        {{- range .Workers }}
        cac-{{.ServerId}}:
          ansible_host: {{.IpAddress}}
          cac_netmask: {{.Netmask}}
          cac_gateway: {{.Gateway}}
          cac_ipv6_address: {{.Ipv6Address}}
          cac_ipv6_gateway: {{.Ipv6Gateway}}
        {{- end}}
`

type Server struct {
	ServerId    int64
	IpAddress   net.IP
	Netmask     net.IP
	Gateway     net.IP
	Ipv6Address net.IPNet
	Ipv6Gateway net.IP
}

var ansibleInventoryCmd = &cobra.Command{
	Use:   "ansible-inventory",
	Short: "create ansible inventory files",
	Run: func(cmd *cobra.Command, args []string) {

		if viper.GetString("username") == "" || viper.GetString("password") == "" {
			fmt.Printf("Set a username and password with --username and --password ")
			fmt.Printf("or by setting the environment variables CAC_USERNAME and CAC_PASSWORD.\n")
			os.Exit(1)
		}

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		c, err := client.New(
			client.WithUsernameAndPassword(
				viper.GetString("username"),
				viper.GetString("password"),
			),
			client.WithUserAgent("cacctl/"+Version),
		)
		if err != nil {
			panic(err)
		}

		response, err := c.List(ctx)
		if err != nil {
			log.Printf("error: %s", err)
			return
		}

		var servers []Server
		for _, server := range response.Servers {

			ipv6Data, err := c.GetIpv6(ctx, server.ServerId)
			if err != nil {
				log.Printf("error: %s", err)
				return
			}

			server := Server{
				ServerId:    server.ServerId,
				IpAddress:   server.IpAddress,
				Netmask:     server.Netmask,
				Gateway:     server.Gateway,
				Ipv6Address: ipv6Data.Ipv6Address,
				Ipv6Gateway: ipv6Data.Ipv6Gateway,
			}
			servers = append(servers, server)
		}

		tmpl, err := template.New("hosts.yaml").Parse(hostsTemplate)
		if err != nil {
			log.Printf("error parsing hosts template: %s", err)
			return
		}

		masters := viper.GetInt("masters")
		if masters > len(servers) {
			masters = len(servers)
		}
		data := struct {
			Masters []Server
			Workers []Server
		}{
			Masters: servers[len(servers)-masters:],
			Workers: servers[:len(servers)-masters],
		}

		err = tmpl.Execute(os.Stdout, data)
		if err != nil {
			log.Printf("error executing host template: %s", err)
		}
	},
}
