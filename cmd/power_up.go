package cmd

import (
	"context"
	"github.com/hef/cacctl/client"
	"github.com/spf13/cobra"
	"log"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(powerUpCmd)
	setupListFlags(powerUpCmd)
}

var powerUpCmd = &cobra.Command{
	Use:   "power_up",
	Short: "power up servers",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupListFlagBindings(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		createClientAndList(ctx, func(c *client.Client, server *client.Server) {
			err := c.PowerCycle(ctx, client.PowerUp, server.VmName, server.ServerId)
			if err != nil {
				log.Printf("error powering up %s: %s", server.ServerName, err)
			} else {
				log.Printf("cac-%d has been powered up", server.ServerId)
			}
		})
	},
}
