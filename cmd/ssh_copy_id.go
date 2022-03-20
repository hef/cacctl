package cmd

import (
	"context"
	"fmt"
	"github.com/hef/cacctl/internal/client"
	"github.com/hef/cacctl/internal/sshx"
	"github.com/pkg/sftp"
	"github.com/spf13/afero/sftpfs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"log"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(sshCopyIdCmd)
	setupListFlags(sshCopyIdCmd)
	sshCopyIdCmd.LocalNonPersistentFlags().StringP("identify-file", "i", "", "Use the identity file")
	viper.BindPFlag("identify-file", sshCopyIdCmd.LocalNonPersistentFlags().Lookup("identify-file"))
}

var sshCopyIdCmd = &cobra.Command{
	Use:   "ssh-copy-id",
	Short: "deploy ssh keys to servers",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupListFlagBindings(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer cancel()

		key, err := sshx.GetPublicKey()
		if err != nil {
			panic(err)
		}

		createClientAndList(ctx, func(c *client.Client, server *client.Server) {
			log.Printf("deploying key to cac-%d at %s", server.ServerId, server.IpAddress)
			deployKey(ctx, server, key)
		})
	},
}

func deployKey(ctx context.Context, server *client.Server, key []byte) {
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
		log.Printf("failed to connect to %d", server.ServerId)
		return
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {
			log.Printf("Failed to close client connection")
		}
	}(client)

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Printf("failed to create sftp client")
		return
	}

	err = sshx.CopyId(ctx, sftpfs.New(sftpClient), key)
	if err != nil {
		log.Printf("error copying id: %s", err)
	}
}
