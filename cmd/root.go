package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var Version = ""
var Commit = ""
var Date = ""

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("username", "", "CAC Username")
	rootCmd.PersistentFlags().String("password", "", "CAC Password")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
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
}
