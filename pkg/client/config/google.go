package config

import "github.com/ArtisanCloud/MediaXCore/utils/object"

type GoogleYouTubeConfig struct {
	*ClientConfig `yaml:",inline"`
	GetOAuthToken func(key string, refresh bool) (token object.HashMap) `yaml:"token;omitempty" json:"token;omitempty"`
	OauthKey      string                                                `yaml:"oauth_key;omitempty" json:"oauth_key;omitempty"`
}

type GoogleBloggerConfig struct {
	*ClientConfig `yaml:",inline"`
}
