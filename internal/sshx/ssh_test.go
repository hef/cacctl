package sshx

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"github.com/spf13/afero"
	"golang.org/x/crypto/ssh"
	"io"
	"testing"
)

func sampleAuthorizedKey() []byte {
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	pubx, err := ssh.NewPublicKey(pub)
	if err != nil {
		panic(err)
	}
	return ssh.MarshalAuthorizedKey(pubx)
}

func TestH(t *testing.T) {

	appFs := afero.NewMemMapFs()
	keyToDeploy := sampleAuthorizedKey()
	err := CopyId(context.TODO(), appFs, keyToDeploy)
	if err != nil {
		panic(err)
	}

	f, err := appFs.Open(".ssh/authorized_keys")
	if err != nil {
		t.Errorf("got unexpected error opening authorized_keys: %s", err)
	}

	raw, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("got unexpected error reading authorized_keys: %s", err)
	}

	deployedKey, _, _, _, _ := ssh.ParseAuthorizedKey(keyToDeploy)
	found, err := checkForKey(deployedKey, raw)
	if err != nil {
		t.Errorf("got unexpected error checking for key: %s", err)
	}

	if !found {
		t.Errorf("expected to find key in authorized_keys but didn't")
	}

	secondKey := sampleAuthorizedKey()
	err = CopyId(context.TODO(), appFs, secondKey)
	if err != nil {
		panic(err)
	}

	f, err = appFs.Open(".ssh/authorized_keys")
	if err != nil {
		t.Errorf("got unexpected error opening authorized_keys: %s", err)
	}

	found, err = checkForKey(deployedKey, raw)
	if err != nil {
		t.Errorf("got unexepcted error checking for first key: %s", err)
	}

	if !found {
		t.Errorf("first key no longer present")
	}

	secondDeployedKey, _, _, _, _ := ssh.ParseAuthorizedKey(secondKey)

	found, err = checkForKey(secondDeployedKey, raw)
	if err != nil {
		t.Errorf("got unexpected error checking for second key: %s", err)
	}
	if !found {
		t.Errorf("second key did not get deployed successfully")
	}

}
