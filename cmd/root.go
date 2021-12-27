package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
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
