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
	"text/tabwriter"
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
		)
		if err != nil {
			panic(err)
		}

		response, err := c.List(ctx)
		if err != nil {
			log.Printf("error: %s", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tIP\tCPU\tRAM\tSSD\tPACKAGE")
		for _, server := range response.Servers {
			fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\t%d\t%s\n", server.ServerId, server.ServerName, server.IpAddress, server.CpuCount, server.RamMB, server.SsdGB, server.Package)
		}
		w.Flush()
	},
}
