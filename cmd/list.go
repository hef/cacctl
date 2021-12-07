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
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Servers",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		c, err := client.New()
		if err != nil {
			panic(err)
		}

		response, err := c.List(ctx)
		if err != nil {
			panic(err)
		}
		log.Printf("list response: %+v", response)

	},
}
