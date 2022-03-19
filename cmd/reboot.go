package cmd

import (
	"context"
	"github.com/hef/cacctl/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

func init() {
	rootCmd.AddCommand(rebootCmd)
}

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "reboot all Servers",
	Run: func(cmd *cobra.Command, args []string) {

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

		serverList, err := c.List(ctx)
		if err != nil {
			log.Printf("Failed to get server list")
			return
		}

		wg := sync.WaitGroup{}
		for _, server := range serverList.Servers {
			wg.Add(1)
			go func(server client.Server) {

				err = c.PowerCycle(ctx, client.Reboot, server.VmName, server.ServerId)
				if err != nil {
					log.Printf("error rebooting %s: %s", server.ServerName, err)
				} else {
					log.Printf("cac-%d has been rebooted", server.ServerId)
				}
				wg.Done()
			}(server)
		}
		wg.Wait()
	},
}
