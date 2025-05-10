// Package config 提供了 MediaX 多平台客户端的配置管理
//
// 本包实现了对各个平台的配置结构定义，包括：
// - 基础配置：API地址、超时设置等
// - OAuth配置：认证信息、Token管理等
// - 平台配置：各平台特定的配置项
//
// 使用说明：
// 1. 首先需要创建对应平台的配置实例
// 2. 填写必要的认证信息和配置参数
// 3. 将配置实例传入对应的客户端创建方法
//
// 详细文档请参考各平台的开发者文档：
// - 微信：https://developers.weixin.qq.com/doc/
// - YouTube：https://developers.google.com/youtube/v3/docs
// - 抖音：https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/
// - 小红书：https://open.xiaohongshu.com/document
package config

import "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"

// MediaXConfig 全局配置结构体
// 包含日志配置等全局设置
type MediaXConfig struct {
	Logger *config.LogConfig // 日志配置，用于设置日志级别、输出方式等
}

// 各平台的 API 接口地址常量定义
const (
	// Wechat 微信公众号接口地址
	WechatAppAPIUrl    string = "https://api.weixin.qq.com"               // 微信公众号 API 基础地址
	WechatAuthTokenUrl string = "https://api.weixin.qq.com/cgi-bin/token" // 微信公众号获取 access_token 地址

	// Google API 接口地址
	GoogleAppAPIUrl     string = "https://www.googleapis.com"         // Google API 基础地址
	GoogleYoutubeAPIUrl string = "https://www.googleapis.com/youtube" // YouTube API 地址
	GoogleBloggerAPIUrl string = "https://www.googleapis.com/blogger" // Blogger API 地址

	// ByteDance 字节跳动接口地址
	ByteDanceDouYinAPIUrl       string = "https://open.douyin.com/"                    // 抖音开放平台 API 地址
	ByteDanceDouYinAuthTokenUrl string = "https://open.douyin.com/oauth/client_token/" // 抖音获取 client_token 地址

	// RedBook 小红书接口地址
	RedBookJuGuangAPIUrl string = "https://adapi.xiaohongshu.com/" // 小红书聚光平台 API 地址
)

// LocalConfig 本地配置结构体
// 包含所有支持平台的配置信息
type LocalConfig struct {
	*WeChatOfficialAccountConfig `yaml:"wechat_official_account_config" json:"wechat_official_account_config"` // 微信公众号配置
	*GoogleYouTubeConfig         `yaml:"google_youtube_config" json:"google_youtube_config"`                   // YouTube 配置
	*GoogleBloggerConfig         `yaml:"google_blogger_config" json:"google_blogger_config"`                   // Blogger 配置
	*ByteDanceDouYinConfig       `yaml:"byte_dance_douyin_config" json:"douyin_config"`                        // 抖音开放平台配置
	*RedBookJuGuangConfig        `yaml:"redbook_juguang_config" json:"redbook_juguang_config"`                 // 小红书聚光平台配置
	*BiliBiliConfig              `yaml:"bilbili_config" json:"bilbili_config"`                                 // B站开放平台配置
}

// BaseConfig 基础配置结构体
// 包含所有平台通用的基础配置项
type BaseConfig struct {
	ApiUrl      string  `yaml:"api_url" json:"api_url"`             // API 基础地址
	ProxyApiUrl string  `yaml:"proxy_api_url" json:"proxy_api_url"` // API 代理地址，用于特殊网络环境
	Timeout     float64 `yaml:"timeout" json:"timeout"`             // 请求超时时间（秒）
	HttpDebug   bool    `yaml:"http_debug" json:"http_debug"`       // 是否开启 HTTP 调试模式
}

// OAuthConfig OAuth2.0 认证配置结构体
// 用于需要 OAuth2.0 认证的平台，如 YouTube、小红书等
// 可以参考 YouTube 的 OAuth 流程：https://developers.google.cn/youtube/v3/guides/authentication?hl=zh-cn
type OAuthConfig struct {
	// OAuth 认证相关 URL
	AccessTokenUrl      string `yaml:"access_token_url,omitempty" json:"access_token_url,omitempty"`             // 获取 access token 的地址
	ProxyAccessTokenUrl string `yaml:"proxy_access_token_url,omitempty" json:"proxy_access_token_url,omitempty"` // 代理获取 access token 的地址
	OAuthUrl            string `yaml:"oauth_url,omitempty" json:"oauth_url,omitempty"`                           // OAuth 授权页面地址
	ProxyOAuthUrl       string `yaml:"proxy_oauth_url,omitempty" json:"proxy_oauth_url,omitempty"`               // OAuth 授权页面代理地址

	// OAuth 客户端认证信息
	ClientID     string `yaml:"client_id" json:"client_id"`                           // OAuth 客户端ID（必填）
	ClientSecret string `yaml:"client_secret" json:"client_secret"`                   // OAuth 客户端密钥（必填）
	RedirectUrl  string `yaml:"redirect_url,omitempty" json:"redirect_url,omitempty"` // 授权回调地址

	// OAuth 授权参数
	Name      string `yaml:"name,omitempty" json:"name,omitempty"`             // 应用名称
	Scope     string `yaml:"scope,omitempty" json:"scope,omitempty"`           // 授权范围
	State     string `yaml:"state,omitempty" json:"state,omitempty"`           // 状态参数，用于防止 CSRF 攻击
	GrantType string `yaml:"grant_type,omitempty" json:"grant_type,omitempty"` // 授权类型，如 authorization_code

	// Token 相关
	RefreshTokenUri string `yaml:"refresh_token_uri,omitempty" json:"refresh_token_uri,omitempty"` // 刷新 token 的地址
	AccessToken     string `yaml:"access_token,omitempty" json:"access_token,omitempty"`           // 访问令牌
	RefreshToken    string `yaml:"refresh_token,omitempty" json:"refresh_token,omitempty"`         // 刷新令牌
	TokenExpiry     int64  `yaml:"token_expiry,omitempty" json:"token_expiry,omitempty"`           // 访问令牌过期时间戳
}

// ClientConfig 客户端配置结构体
// 组合了基础配置和 OAuth 配置
type ClientConfig struct {
	*BaseConfig  `yaml:",inline"` // 基础配置
	*OAuthConfig `yaml:"oauth"`   // OAuth 配置
}
