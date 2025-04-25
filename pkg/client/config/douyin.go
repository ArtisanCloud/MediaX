package config

type ByteDanceDouYinConfig struct {
	*ClientConfig `yaml:",inline"`
	//GetOAuthToken func(key string, refresh bool) (token object.HashMap) `yaml:"token;omitempty" json:"token;omitempty"`
}
