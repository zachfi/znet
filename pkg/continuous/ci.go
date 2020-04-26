package continuous

import (
	git "github.com/go-git/go-git/v5"
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

	return nil
}
