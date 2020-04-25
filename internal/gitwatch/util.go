package gitwatch

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
)

func CacheRepo(url string, cacheDir string, sshPublicKey string) error {
	_, err := os.Stat(cacheDir)
	if err != nil {
		cloneOpts := &git.CloneOptions{
			URL:      url,
			Progress: nil,
			// Progress: os.Stdout,
		}

		publicKey, err := SSHPublicKey(url, sshPublicKey)
		if err != nil {
			return err
		}

		if publicKey != nil {
			cloneOpts.Auth = publicKey
		}

		log.Infof("cloning repo from origin %s", url)
		_, cloneErr := git.PlainClone(cacheDir, false, cloneOpts)
		if cloneErr != nil {
			return fmt.Errorf("error while cloning %s: %s", url, cloneErr)
		}
	}

	return nil
}

func SSHPublicKey(url string, SSHKeyPath string) (*ssh.PublicKeys, error) {

	// For URLs that don't start with http and when a SSHKeyPath is set, we
	// shold load the ssh key to proceed.
	if !strings.HasPrefix(url, "https") && SSHKeyPath != "" {
		var publicKey *ssh.PublicKeys
		sshKey, _ := ioutil.ReadFile(SSHKeyPath)
		publicKey, keyError := ssh.NewPublicKeys("git", sshKey, "")
		if keyError != nil {
			return nil, fmt.Errorf("error while loading public key: %s", keyError)
		}

		return publicKey, nil
	}

	return nil, nil
}
