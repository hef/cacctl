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
	rootCmd.AddCommand(rebootCmd)
	setupListFlags(rebootCmd)
}

var powerDownCmd = &cobra.Command{
	Use:   "power_down",
	Short: "power down servers",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupListFlagBindings(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		createClientAndList(ctx, func(c *client.Client, server *client.Server) {
			err := c.PowerCycle(ctx, client.PowerDown, server.VmName, server.ServerId)
			if err != nil {
				log.Printf("error powering down %s: %s", server.ServerName, err)
			} else {
				log.Printf("cac-%d has been rebooted", server.ServerId)
			}
		})
	},
}
