package config

type GoogleYouTubeConfig struct {
	*ClientConfig `yaml:",inline"`
}

type GoogleBloggerConfig struct {
	*ClientConfig `yaml:",inline"`
}
