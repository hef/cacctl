package cmd

import (
	"context"
	"github.com/hef/cacctl/internal/client"
	"github.com/spf13/cobra"
	"log"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	setupListFlags(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a server",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupListFlagBindings(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		createClientAndList(ctx, func(c *client.Client, server *client.Server) {
			err := c.Delete(ctx, server.ServerId, server.CustomerId, server.ServerName, false)
			if err != nil {
				log.Printf("error deleting server: %s", err)
			}
		})
	},
}
