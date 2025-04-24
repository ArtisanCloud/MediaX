package config

import "github.com/ArtisanCloud/MediaXCore/utils/object"

type RedBookJuGuangConfig struct {
	*ClientConfig `yaml:",inline"`
	GetOAuthToken func(key string, refresh bool) (token object.HashMap) `yaml:"token;omitempty" json:"token;omitempty"`
}
