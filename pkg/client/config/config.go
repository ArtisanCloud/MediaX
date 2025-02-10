package config

import "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"

type MediaXConfig struct {
	Logger *config.LogConfig
}

type LocalConfig struct {
	*WeChatOfficialAccountConfig `yaml:"wechat_config" json:"wechat_config"`
	*DouYinConfig                `yaml:"douyin_config" json:"douyin_config"`
	*RedBookConfig               `yaml:"redbook_config" json:"redbook_config"`
}

type AppConfig struct {
	BaseUri   string  `yaml:"base_uri" json:"base_uri"`
	ProxyUri  string  `yaml:"proxy_uri" json:"proxy_uri"`
	Timeout   float64 `yaml:"timeout" json:"timeout"`
	AppID     string  `yaml:"app_id" json:"app_id"`
	AppSecret string  `yaml:"app_secret" json:"app_secret"`
	HttpDebug bool    `yaml:"http_debug" json:"http_debug"`
}
