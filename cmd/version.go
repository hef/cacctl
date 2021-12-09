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
		if Version != "" {
			fmt.Println("Version: " + Version)
		}
		if Commit != "" {
			fmt.Println("Commit: " + Commit)
		}
		if Date != "" {
			fmt.Println("Build Date: " + Date)
		}
		fmt.Println(runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH)
	},
}
