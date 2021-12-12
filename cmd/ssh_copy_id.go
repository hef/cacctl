package cmd

import (
	"context"
	"fmt"
	"github.com/hef/cacctl/internal/client"
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(sshCopyIdCmd)
	sshCopyIdCmd.PersistentFlags().StringP("identify-file", "i", "", "Use the identity file")
	viper.BindPFlag("cpu", sshCopyIdCmd.PersistentFlags().Lookup("cpu"))
}

var sshCopyIdCmd = &cobra.Command{
	Use:   "ssh-copy-id",
	Short: "deploy ssh keys to servers",
	Run: func(cmd *cobra.Command, args []string) {

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

		serverList, err := c.List(ctx)
		if err != nil {
			log.Printf("Failed to get server list")
		}

		for _, server := range serverList.Servers {
			err = sshCopyId(ctx, &server)
		}
	},
}

func sshCopyId(ctx context.Context, server *client.Server) error {
	ip := server.IpAddress
	config := ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip.String(), 22), &config)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}

	f, err := sftpClient.Open(".profile")
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, f)
	return nil
}
