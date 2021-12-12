package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/hef/cacctl/internal/client"
	"github.com/mitchellh/go-homedir"
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

		key, err := getPublicKey()
		if err != nil {
			panic(err)
		}

		serverList, err := c.List(ctx)
		if err != nil {
			log.Printf("Failed to get server list")
		}

		for _, server := range serverList.Servers {
			err = sshCopyId(ctx, &server, key)
			if err != nil {
				log.Printf("error copying id: %s", err)
			}
		}
	},
}

func sshCopyId(ctx context.Context, server *client.Server, key []byte) error {
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

	f, err := sftpClient.Open(".ssh/authorized_keys")
	if errors.Is(err, os.ErrNotExist) {
		f, err := sftpClient.Create(".ssh/authorized_keys")
		if errors.Is(err, os.ErrNotExist) {
			err = sftpClient.Mkdir(".ssh")
			if err != nil {
				return err
			}
			f, err = sftpClient.Create(".ssh/authorized_keys")
		}
		if err != nil {
			return err
		}
		defer f.Close()
		io.Copy(f, bytes.NewReader(key))
		return nil
	}
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, f)
	return nil
}

func getPublicKey() ([]byte, error) {
	identityFile := viper.GetString("identity-file")
	if identityFile == "" {
		keyPaths := []string{
			"~/.ssh/id_ed25519.pub",
			"~/.ssh/id_dsa.pub",
			"~/.ssh/id_rsa.pub",
		}
		for _, keyPath := range keyPaths {
			path, err := homedir.Expand(keyPath)
			if err != nil {
				return nil, err
			}
			_, err = os.Stat(path)
			if err == nil {
				identityFile = path
				break
			}
		}
	}

	if len(identityFile) < 4 || identityFile[len(identityFile)-4:] != ".pub" {
		identityFile = identityFile + ".pub"
	}

	f, err := os.Open(identityFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
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
