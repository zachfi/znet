package continuous

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/kballard/go-shellquote"
	log "github.com/sirupsen/logrus"
)

// CacheRepo is used to clone a repo using SSH publich key authentication.
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

// SSHPublicKey is used to load an SSH public key for use in authenticating to a git server.
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

// FetchRemote performs a fetch of a given repository, returning a the new branch heads and new tags.
func FetchRemote(repo *git.Repository, sshPublicKey *ssh.PublicKeys) (map[string]string, []string, error) {
	newHeads := make(map[string]string)
	newTags := make([]string, 0)
	var err error

	beforeHeads, beforeTags, err := RepoRefs(repo)
	if err != nil {
		return newHeads, newTags, err
	}

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

	afterHeads, afterTags, err := RepoRefs(repo)
	if err != nil {
		return newHeads, newTags, err
	}

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

// RepoRefs is used to retrieve the references of a working repo.  This is
// useful to check the state before and after a fetch, to determine the
// difference.
func RepoRefs(repo *git.Repository) (map[string]string, map[string]string, error) {
	heads := make(map[string]string)
	tags := make(map[string]string)

	if repo == nil {
		return heads, tags, fmt.Errorf("unable to operate on nil repository")
	}

	refs, err := repo.References()
	if err != nil {
		return heads, tags, err
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
		return heads, tags, err
	}

	return heads, tags, nil
}

func Build(commandLine string, cacheDir string) (BuildResult, error) {
	var ev BuildResult
	var args []string
	var err error

	parts := strings.SplitN(commandLine, " ", 2)

	commandName := parts[0]

	if len(parts) > 0 {
		args, err = shellquote.Split(parts[1])
		if err != nil {
			return ev, err
		}
	}

	cmd := exec.Command(commandName, args...)
	cmd.Dir = cacheDir

	log.Debugf("executing command: %+v", *cmd)
	startTime := time.Now()
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("command execution failed: %s", err)
	}

	duration := time.Since(startTime)

	log.Debugf("command output: %+v", string(output))

	now := time.Now()

	ev = BuildResult{
		Time:     &now,
		Command:  commandName,
		Args:     args,
		Dir:      cacheDir,
		Output:   output,
		ExitCode: cmd.ProcessState.ExitCode(),
		Duration: duration,
	}

	return ev, nil
}
