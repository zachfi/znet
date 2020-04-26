package builder

type Config struct {
	CacheDir   string `yaml:"cache_dir"`
	SSHKeyPath string `yaml:"ssh_key_path"`
}
