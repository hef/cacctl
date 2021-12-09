package cmd

import (
	"context"
	"github.com/hef/cacctl/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().String("username", "", "cac username")
	listCmd.PersistentFlags().String("password", "", "cac password")
	viper.BindPFlag("username", listCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", listCmd.PersistentFlags().Lookup("password"))
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Servers",
	Run: func(cmd *cobra.Command, args []string) {

		if viper.GetString("username") == "" || viper.GetString("password") == "" {
			log.Printf("Set a username and password with --username and --password")
			log.Printf("or by setting the environment variables CAC_USERNAME and CAC_PASSWORD")
			os.Exit(1)
		}

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		c, err := client.New(
			client.WithUsernameAndPassword(
				viper.GetString("username"),
				viper.GetString("password"),
			),
		)
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
