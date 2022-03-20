package cmd

import (
	"context"
	"fmt"
	"github.com/hef/cacctl/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"sync"
)

var Version = ""
var Commit = ""
var Date = ""

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("username", "", "CAC Username")
	rootCmd.PersistentFlags().String("password", "", "CAC Password")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

var rootCmd = &cobra.Command{
	Use: "cacctl",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetEnvPrefix("CAC")
	viper.AutomaticEnv()
	configPath := viper.GetString("config")
	if configPath == "" {
		configDir, err := os.UserConfigDir()
		if err == nil {
			viper.AddConfigPath(path.Join(configDir, "cacctl"))
		}
		err = viper.ReadInConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && err != nil {
			log.Printf("error reading config file %s: %s", viper.ConfigFileUsed(), err)
		}
	} else {
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Printf("error reading config file %s: %s", viper.ConfigFileUsed(), err)
		}
	}
}

func createClientAndList(ctx context.Context, lambda func(c *client.Client, server *client.Server)) (c *client.Client, r *client.ListResponse, err error) {

	c, err = client.New(
		client.WithUsernameAndPassword(
			viper.GetString("username"),
			viper.GetString("password"),
		),
		client.WithUserAgent("cacctl/"+Version),
	)

	if err != nil {
		return nil, nil, err
	}

	search := viper.GetString("search")
	limit := viper.GetInt("limit")
	filter := viper.GetString("filter")
	if search == "" && limit == 25 && filter == "All" {
		r, err = c.List(ctx)
	} else {
		r, err = c.ListWithFilter(ctx, search, limit, client.ListFilterFromString(filter))
	}

	if lambda != nil {
		wg := sync.WaitGroup{}
		for _, server := range r.Servers {
			wg.Add(1)
			go func(server client.Server) {
				lambda(c, &server)
				wg.Done()
			}(server)
		}
		wg.Wait()
	}
	return
}
