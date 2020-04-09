package gitwatch

type Config struct {
	CacheDir string `yaml:"cache_dir"`
	Repos    []Repo `yaml:"repos"`
}

type Repo struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}
