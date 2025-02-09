package config

import "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"

type MediaXConfig struct {
	Logger *config.LogConfig
}

type LocalConfig struct {
	*WeChatConfig  `yaml:"wechat_config" json:"weChatConfig"`
	*DouYinConfig  `yaml:"douyin_config" json:"douYinConfig"`
	*RedBookConfig `yaml:"redbook_config" json:"redBookConfig"`
}

type AppConfig struct {
	BaseUri   string  `yaml:"base_uri" json:"baseUri"`
	ProxyUri  string  `yaml:"proxy_uri" json:"proxyUri"`
	Timeout   float64 `yaml:"timeout" json:"timeout"`
	AppID     string  `yaml:"app_id" json:"appId"`
	AppSecret string  `yaml:"app_secret" json:"appSecret"`
	HttpDebug bool    `yaml:"http_debug" json:"httpDebug"`
}
