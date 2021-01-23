package config

type BuilderConfig struct {
	CacheDir   string `yaml:"cache_dir"`
	SSHKeyPath string `yaml:"ssh_key_path"`
}

type RepoConfig struct {
	OnTag    []string `yaml:"on_tag"`
	OnCommit []string `yaml:"on_commit"`
	Branches []string `yaml:"branches"`
}
