package cmd

import (
	"context"
	"fmt"
	"github.com/hef/cacctl/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().Int("cpu", 0, "Number of CPUs [1 - 10]")
	buildCmd.Flags().Int("ram", 0, "Ram in MB, valid values are [512, 1024, 1536, 2048, 2560, 3072, 4096, 5120, 6144, 7168, 8192]")
	buildCmd.Flags().Int("storage", 0, "Storage to allocate, in GB, valid values are multiples of 10")
	buildCmd.Flags().Bool("ha", false, "Enable High Availability [true|false]")
	buildCmd.Flags().Bool("encryption", false, "Enable Encryption [true|false]")
	buildCmd.Flags().String("os", "", `Operating System ["CentOS 7.9 64bit", "CentOS 8.3 64bit", "Debian 9.13 64Bit", "FreeBSD 12.2 64bit", "Ubuntu 18.04 LTS 64bit"]`)
	viper.BindPFlag("cpu", buildCmd.Flags().Lookup("cpu"))
	viper.BindPFlag("ram", buildCmd.Flags().Lookup("ram"))
	viper.BindPFlag("storage", buildCmd.Flags().Lookup("storage"))
	viper.BindPFlag("ha", buildCmd.Flags().Lookup("ha"))
	viper.BindPFlag("encryption", buildCmd.Flags().Lookup("encryption"))
	viper.BindPFlag("os", buildCmd.Flags().Lookup("os"))
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

		request := client.BuildRequest{
			Cpu:              viper.GetInt("cpu"),
			Ram:              viper.GetInt("ram"),
			Storage:          viper.GetInt("storage"),
			HighAvailability: viper.GetBool("ha"),
			Encryption:       viper.GetBool("encryption"),
			OS:               viper.GetString("os"),
		}

		_, err = c.Build(ctx, &request)
		if err != nil {
			log.Printf("error building: %s", err)
		}
	},
}
