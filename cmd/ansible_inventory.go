//go:build dev

package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(ansibleInventoryCmd)
}

var ansibleInventoryCmd = &cobra.Command{
	Use:   "ansible-inventory",
	Short: "create ansible inventory files",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
