package cmd

import (
	"context"
	"github.com/hef/cacctl/internal/client"
	"github.com/spf13/cobra"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

func init() {
	rootCmd.AddCommand(rebootCmd)
	setupListFlags(rebootCmd)
}

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "reboot servers",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		c, response, err := createClientAndList(ctx)
		if err != nil {
			log.Printf("error creating client: %s", err)
			return
		}

		wg := sync.WaitGroup{}
		for _, server := range response.Servers {
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
