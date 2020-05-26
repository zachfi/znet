package continuous

import (
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

type CI struct {
	URL        string
	CacheDir   string
	SSHKeyPath string
}

func NewCI(url, cacheDir, sshKeyPath string) *CI {
	c := &CI{
		URL:        url,
		CacheDir:   cacheDir,
		SSHKeyPath: sshKeyPath,
	}

	return c
}

func (c *CI) LatestTag() string {
	err := CacheRepo(c.URL, c.CacheDir, c.SSHKeyPath)
	if err != nil {
		log.Errorf("error while caching repo %s: %s", c.URL, err)
	}

	r, err := git.PlainOpen(c.CacheDir)
	if err != nil {
		log.Error(err)
	}

	_, tags, err := RepoRefs(r)
	if err != nil {
		log.Error(err)
	}

	var lastTag string
	var prefix string

	for n := range tags {
		if lastTag == "" {
			lastTag = n
		}

		v, err := semver.ParseTolerant(n)
		if err != nil {
			log.Error(err)
			continue
		}

		nv, err := semver.ParseTolerant(lastTag)
		if err != nil {
			log.Error(err)
			continue
		}

		if v.GT(nv) {
			lastTag = v.String()

			if strings.HasPrefix(n, "v") {
				prefix = "v"
			} else {
				prefix = ""
			}
		}
	}

	return fmt.Sprintf("%s%s", prefix, lastTag)
}

// Fetch performs a git-fetch from the repo origin.
func (c *CI) Fetch() (map[string]string, []string, error) {

	err := CacheRepo(c.URL, c.CacheDir, c.SSHKeyPath)
	if err != nil {
		log.Errorf("error while caching repo %s: %s", c.URL, err)
	}

	r, err := git.PlainOpen(c.CacheDir)
	if err != nil {
		log.Error(err)
	}

	publicKey, err := SSHPublicKey(c.URL, c.SSHKeyPath)
	if err != nil {
		log.Error(err)
	}

	return FetchRemote(r, publicKey)
}

func (c *CI) CheckoutTag(tag string) error {
	r, err := git.PlainOpen(c.CacheDir)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	ref, err := r.Tag(tag)
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{Hash: ref.Hash()})
	if err != nil {
		return err
	}

	err = w.Reset(&git.ResetOptions{Mode: git.HardReset})
	if err != nil {
		return err
	}

	return nil
}

func (c *CI) CheckoutHash(commit string) error {
	r, err := git.PlainOpen(c.CacheDir)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	hash := plumbing.NewHash(commit)

	err = w.Checkout(&git.CheckoutOptions{Hash: hash})
	if err != nil {
		return err
	}

	return nil
}
