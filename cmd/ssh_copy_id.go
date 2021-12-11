//go:build dev

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(sshCopyIdCmd)
	sshCopyIdCmd.PersistentFlags().StringP("identify-file", "i", "", "Use the identity file")
	viper.BindPFlag("cpu", sshCopyIdCmd.PersistentFlags().Lookup("cpu"))
}

var sshCopyIdCmd = &cobra.Command{
	Use:   "ssh-copy-id",
	Short: "List all Servers",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
