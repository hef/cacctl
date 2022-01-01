package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"runtime"
	"syscall"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configure credentials",
	Run: func(cmd *cobra.Command, args []string) {

		configPath := ""
		if viper.ConfigFileUsed() != "" {
			configPath = viper.ConfigFileUsed()
		} else {
			configDir, err := os.UserConfigDir()
			if err == nil {
				configPath = path.Join(configDir, "cacctl", "config.yaml")
			}
		}
		fmt.Printf("using %s\n", configPath)

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("CAC Username: ")
		username, err := reader.ReadString('\n')
		if runtime.GOOS == "windows" {
			username = username[:len(username)-2]
		} else {
			username = username[:len(username)-1]
		}
		if err != nil {
			panic(err)
		}

		fmt.Print("CAC Password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		fmt.Print("\n")

		config := struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}{
			username,
			string(password),
		}

		_ = os.MkdirAll(path.Dir(viper.ConfigFileUsed()), os.ModePerm)
		f, err := os.Create(configPath)
		err = yaml.NewEncoder(f).Encode(&config)
		if err != nil {
			panic(err)
		}

		if err != nil {
			panic(err)
		}
	},
}
