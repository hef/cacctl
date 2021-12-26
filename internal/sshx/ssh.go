package sshx

import (
	"bytes"
	"context"
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
)

func CopyId(ctx context.Context, fs afero.Fs, key []byte) error {

	f, err := fs.Open(".ssh/authorized_keys")
	if errors.Is(err, os.ErrNotExist) {

		f, err := fs.Create(".ssh/authorized_keys")
		if errors.Is(err, os.ErrNotExist) {
			err = fs.Mkdir(".ssh", 0644)
			if err != nil {
				return err
			}
			f, err = fs.Create(".ssh/authorized_keys")
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

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	parsedKey, _, _, _, err := ssh.ParseAuthorizedKey(key)
	if err != nil {
		return err
	}
	keyExists, err := checkForKey(parsedKey, raw)
	if err != nil {
		return err
	}

	if keyExists {
		return nil
	}

	newData := append([]byte{'\n'}, key...)
	_, err = f.Write(newData)
	if err != nil {
		return err
	}
	return f.Close()
}

func merge(existingFile io.Reader, keydata []byte) (x []byte, err error) {

	raw, err := ioutil.ReadAll(existingFile)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePublicKey(keydata)
	if err != nil {
		return nil, err
	}

	present, err := checkForKey(key, raw)
	if err != nil {
		return nil, err
	}
	if present {
		return nil, nil
	}

	x = append(raw, keydata...)
	return x, nil
}

func GetPublicKey() ([]byte, error) {
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

func checkForKey(deployedKey ssh.PublicKey, authorizedKeys []byte) (bool, error) {
	var extractedKey ssh.PublicKey
	var err error
	for rest := authorizedKeys; rest != nil; extractedKey, _, _, rest, err = ssh.ParseAuthorizedKey(rest) {
		if err != nil {
			return false, err
		}
		if extractedKey != nil && bytes.Compare(extractedKey.Marshal(), deployedKey.Marshal()) == 0 {
			return true, nil
		}
	}
	return false, nil
}
