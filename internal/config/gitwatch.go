package config

type GitWatchConfig struct {
	CacheDir    string               `yaml:"cache_dir"`
	Repos       []GitWatchRepo       `yaml:"repos"`
	Interval    int                  `yaml:"interval"`
	SSHKeyPath  string               `yaml:"ssh_key_path"`
	Collections []GitWatchCollection `yaml:"collections"`
}

type GitWatchCollection struct {
	Name     string
	Repos    []GitWatchRepo `yaml:"repos"`
	Interval int            `yaml:"interval"`
}

type GitWatchRepo struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}
