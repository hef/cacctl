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
	rootCmd.AddCommand(deleteCmd)
	setupListFlags(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "build a server",
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
				err = c.Delete(ctx, server.ServerId, server.CustomerId, server.ServerName, false)
				wg.Done()
			}(server)
		}
		wg.Wait()
	},
}
