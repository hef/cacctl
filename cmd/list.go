package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
)

func setupListFlags(command *cobra.Command) {
	command.Flags().String("search", "", `search string`)
	command.Flags().Int("limit", 25, `[5, 10, 25, 50, 100, 150, 200]`)
	command.Flags().String("filter", "All", `[All, PoweredOn, PoweredOff, Installing, Pending, Installed]`)

}

func setupListFlagBindings(command *cobra.Command) {
	viper.BindPFlag("search", command.Flags().Lookup("search"))
	viper.BindPFlag("limit", command.Flags().Lookup("limit"))
	viper.BindPFlag("filter", command.Flags().Lookup("filter"))
}

func init() {
	rootCmd.AddCommand(listCmd)
	setupListFlags(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Servers",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupListFlagBindings(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		_, response, err := createClientAndList(ctx, nil)
		if err != nil {
			log.Printf("error creating client: %s", err)
			return
		}

		if response != nil {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tSTATUS\tIP\tCPU\tRAM\tSSD\tOS\tPACKAGE")
			for _, server := range response.Servers {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%d\t%d\t%s\t%s\n", server.ServerId, server.ServerName, server.Status, server.IpAddress, server.CpuCount, server.RamMB, server.SsdGB, server.CurrentOs, server.Package)
			}
			w.Flush()
		}
	},
}
