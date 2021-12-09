package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: " + Version)
		fmt.Println("Commit: " + Commit)
		fmt.Println("Build Date: " + Date)
		fmt.Println("Go Version: " + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH)
	},
}
