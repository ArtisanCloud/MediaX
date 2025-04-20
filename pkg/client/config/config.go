package config

import "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"

type MediaXConfig struct {
	Logger *config.LogConfig
}

type LocalConfig struct {
	*WeChatOfficialAccountConfig `yaml:"wechat_official_account_config" json:"wechat_official_account_config"`
	*GoogleYouTubeConfig         `yaml:"google_youtube_config" json:"google_youtube_config"`
	*GoogleBloggerConfig         `yaml:"google_blogger_config" json:"google_blogger_config"`
	*DouYinConfig                `yaml:"douyin_config" json:"douyin_config"`
	*RedBookConfig               `yaml:"redbook_config" json:"redbook_config"`
}
type BaseConfig struct {
	ApiUrl      string  `yaml:"api_url" json:"api_url"`
	ProxyApiUrl string  `yaml:"proxy_api_url" json:"proxy_api_url"`
	Timeout     float64 `yaml:"timeout" json:"timeout"`
	HttpDebug   bool    `yaml:"http_debug" json:"http_debug"`
}

//type AppConfig struct {
//	ApiUrl      string `yaml:"api_url" json:"api_url"`
//	ProxyApiUrl string `yaml:"proxy_api_url" json:"proxy_api_url"`
//
//	AppID     string `yaml:"app_id" json:"app_id"`
//	AppSecret string `yaml:"app_secret" json:"app_secret"`
//}

type OAuthConfig struct {
	// 可以参考youtube的oauth流程：https://developers.google.cn/youtube/v3/guides/authentication?hl=zh-cn
	// 必填项
	OAuthUrl       string `yaml:"oauth_url,omitempty" json:"oauth_url,omitempty"`
	ProxyOAuthUrl  string `yaml:"proxy_oauth_url,omitempty" json:"proxy_oauth_url,omitempty"`
	ClientID       string `yaml:"client_id" json:"client_id"`                                   // OAuth 客户端ID (必填)
	ClientSecret   string `yaml:"client_secret" json:"client_secret"`                           // OAuth 客户端密钥 (必填)
	RedirectUrl    string `yaml:"redirect_url,omitempty" json:"redirect_url,omitempty"`         // 重定向 URL (可选)
	AccessTokenUrl string `yaml:"access_token_url,omitempty" json:"access_token_url,omitempty"` // 获取 access token 的 URI (可选)

	// 可选项，使用 omitempty 来省略空值字段
	Name            string `yaml:"name,omitempty" json:"name,omitempty"`                           // 授权范围 (可选，只有在存在时才包括)
	Scope           string `yaml:"scope,omitempty" json:"scope,omitempty"`                         // 授权范围 (可选，只有在存在时才包括)
	State           string `yaml:"state,omitempty" json:"state,omitempty"`                         // 状态参数 (可选，只有在存在时才包括)
	GrantType       string `yaml:"grant_type,omitempty" json:"grant_type,omitempty"`               // 授权类型 (可选，只有在存在时才包括)
	RefreshTokenUri string `yaml:"refresh_token_uri,omitempty" json:"refresh_token_uri,omitempty"` // 刷新 token 的 URI (可选，只有在存在时才包括)

	// 存储用字段，只有在获取到 token 后才需要
	AccessToken  string `yaml:"access_token,omitempty" json:"access_token,omitempty"`   // 存储获取到的 access token (可选，只有在存在时才包括)
	RefreshToken string `yaml:"refresh_token,omitempty" json:"refresh_token,omitempty"` // 存储获取到的 refresh token (可选，只有在存在时才包括)
	TokenExpiry  int64  `yaml:"token_expiry,omitempty" json:"token_expiry,omitempty"`   // 访问令牌的过期时间 (可选，只有在存在时才包括)
}

type ClientConfig struct {
	*BaseConfig  `yaml:",inline"`
	*OAuthConfig `yaml:"oauth"`
	//*AppConfig   `yaml:"app"`
}
