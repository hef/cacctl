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
	sshCopyIdCmd.PersistentFlags().StringP("identify-file", "i", "", "Use the identity file")
	viper.BindPFlag("identify-file", sshCopyIdCmd.PersistentFlags().Lookup("identify-file"))
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

		key, err := sshx.GetPublicKey()
		if err != nil {
			panic(err)
		}

		serverList, err := c.List(ctx)
		if err != nil {
			log.Printf("Failed to get server list")
			return
		}

		for _, server := range serverList.Servers {

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
				break
			}
			defer client.Close()

			sftpClient, err := sftp.NewClient(client)
			if err != nil {
				log.Printf("failed to create sftp client")
				break
			}

			err = sshx.CopyId(ctx, sftpfs.New(sftpClient), key)
			if err != nil {
				log.Printf("error copying id: %s", err)
			}
		}
	},
}

func mergeKey(authorizedKeys, newKey []byte) (newFile []byte, err error) {

	key, err := ssh.ParsePublicKey(newKey)
	if err != nil {
		return nil, err
	}

	for parseKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeys); rest != nil {
		if err != nil {
			return nil, err
		}
		if bytes.Equal(key.Marshal(), parseKey.Marshal()) {
			return nil, nil
		}

	}

}
