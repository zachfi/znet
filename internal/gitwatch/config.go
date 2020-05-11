package gitwatch

type Config struct {
	CacheDir    string       `yaml:"cache_dir"`
	Repos       []Repo       `yaml:"repos"`
	Interval    int          `yaml:"interval"`
	SSHKeyPath  string       `yaml:"ssh_key_path"`
	Collections []Collection `yaml:"collections"`
}

type Collection struct {
	Name  string
	Repos []Repo `yaml:"repos"`
}

type Repo struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}
