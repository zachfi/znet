package gitwatch

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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

func FetchRemote(repo *git.Repository, sshPublicKey *ssh.PublicKeys) (map[string]string, []string, error) {
	newHeads := make(map[string]string)
	newTags := make([]string, 0)
	var err error

	beforeHeads, beforeTags := RepoRefs(repo)

	fetchOpts := &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []gitConfig.RefSpec{
			"+refs/heads/*:refs/remotes/origin/*",
			"+refs/remotes/*:refs/remotes/origin/*",
		},
		Tags:  git.AllTags,
		Force: true,
	}

	if sshPublicKey != nil {
		fetchOpts.Auth = sshPublicKey
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		log.Errorf("error repo.Remote() %s", err)
	}

	err = remote.Fetch(fetchOpts)
	if err != nil {
		if err.Error() != git.NoErrAlreadyUpToDate.Error() {
			log.Errorf("failed to fetch %s: %s", remote.Config().Name, err)
		} else {
			err = nil
		}
	}

	afterHeads, afterTags := RepoRefs(repo)

	nameMatch := func(refs map[string]string, shortName string) bool {
		for k := range refs {
			if k == shortName {
				return true
			}
		}

		return false
	}

	refMatch := func(refs map[string]string, shortName string, hash string) bool {
		for k, v := range refs {
			if k == shortName {
				if v == hash {
					return true
				}
			}
		}

		return false
	}

	// detect new commits on all branches
	for shortName, hash := range afterHeads {
		//detect new branches
		if !nameMatch(beforeHeads, shortName) {
			newHeads[shortName] = hash
			continue
		}

		// when before did not have this branch
		if !refMatch(beforeHeads, shortName, hash) {
			newHeads[shortName] = hash
		}
	}

	// detect new tags
	for shortName := range afterTags {
		if !nameMatch(beforeTags, shortName) {
			newTags = append(newTags, shortName)
		}
	}

	return newHeads, newTags, err
}

func RepoRefs(repo *git.Repository) (map[string]string, map[string]string) {
	heads := make(map[string]string)
	tags := make(map[string]string)

	refs, err := repo.References()
	if err != nil {
		log.Error(err)
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		// The HEAD is omitted in a `git show-ref` so we ignore the symbolic
		// references, the HEAD
		if ref.Type() == plumbing.SymbolicReference {
			return nil
		}

		// Only inspect the remote references
		if ref.Name().IsRemote() {
			heads[ref.Name().Short()] = ref.Hash().String()
		}

		if ref.Name().IsTag() {
			tags[ref.Name().Short()] = ref.Hash().String()
		}

		return nil
	})

	if err != nil {
		log.Error(err)
	}

	return heads, tags
}
