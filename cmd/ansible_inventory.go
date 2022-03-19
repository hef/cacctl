//go:build dev

package cmd

import (
	"context"
	"fmt"
	"github.com/hef/cacctl/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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
        {{- end}}
    worker:
      hosts:
        {{- range .Workers }}
        cac-{{.ServerId}}:
          ansible_host: {{.IpAddress}}
        {{- end}}
`

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

		tmpl, err := template.New("hosts.yaml").Parse(hostsTemplate)
		if err != nil {
			log.Printf("error parsing hosts template: %s", err)
			return
		}

		masters := viper.GetInt("masters")
		if masters > len(response.Servers) {
			masters = len(response.Servers)
		}
		data := struct {
			Masters []client.Server
			Workers []client.Server
		}{
			Masters: response.Servers[:masters],
			Workers: response.Servers[masters:],
		}

		err = tmpl.Execute(os.Stdout, data)
		if err != nil {
			log.Printf("error executing host template: %s", err)
		}

	},
}
