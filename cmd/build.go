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
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build a server",
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
		resp, err := c.Build(ctx, &client.BuildRequest{})
		log.Printf("+%v", resp)
	},
}
