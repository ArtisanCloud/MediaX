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

import (
	"fmt"
	"strings"

	"github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
)

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

	// BiliBili 哔哩哔哩接口地址
	BiliBiliAPIUrl string = "https://member.bilibili.com/" // 哔哩哔哩开放平台 API 地址
)

// LocalConfig 本地配置结构体
// 包含所有支持平台的配置信息
type LocalConfig struct {
	AccessTokenProviders  *AccessTokenProvidersConfig  `yaml:"access_token_providers" json:"access_token_providers"`
	ClientTokenProviders  *ClientTokenProvidersConfig  `yaml:"client_token_providers" json:"client_token_providers"`
	SessionTokenProviders *SessionTokenProvidersConfig `yaml:"session_token_providers" json:"session_token_providers"`
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

// AccessTokenProvidersConfig 表示 AccessToken 调试相关的 Provider/App 列表
type AccessTokenProvidersConfig struct {
	Redis     *AccessTokenRedisConfig `yaml:"redis,omitempty" json:"redis,omitempty"`
	Providers []*AccessTokenProvider  `yaml:"providers" json:"providers"`
}

// AccessTokenRedisConfig 定义调试服务使用的 Redis 连接
type AccessTokenRedisConfig struct {
	Addr     string `yaml:"addr,omitempty" json:"addr,omitempty"`
	DB       int    `yaml:"db,omitempty" json:"db,omitempty"`
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

// AccessTokenProvider 代表 Provider 分组（例如 Google、字节跳动）
type AccessTokenProvider struct {
	Code string                    `yaml:"code" json:"code"`
	Name string                    `yaml:"name" json:"name"`
	Apps []*AccessTokenProviderApp `yaml:"apps" json:"apps"`
}

// AccessTokenProviderApp 代表某个 Provider 下的具体 App（例如 YouTube、Blogger）
type AccessTokenProviderApp struct {
	Code         string                 `yaml:"code" json:"code"`
	Name         string                 `yaml:"name" json:"name"`
	ProviderCode string                 `yaml:"provider_code" json:"provider_code"`
	ApiVersion   string                 `yaml:"api_version,omitempty" json:"api_version,omitempty"`
	AuthModes    []*AccessTokenAuthMode `yaml:"auth_modes" json:"auth_modes"`
}

// AccessTokenAuthMode 描述 App 的授权模式（不同 oauth_key / 环境）
type AccessTokenAuthMode struct {
	Key                   string                 `yaml:"key" json:"key"`
	Label                 string                 `yaml:"label,omitempty" json:"label,omitempty"`
	ProviderCode          string                 `yaml:"provider_code,omitempty" json:"provider_code,omitempty"`
	GoogleYouTubeConfig   *GoogleYouTubeConfig   `yaml:"google_youtube_config,omitempty" json:"google_youtube_config,omitempty"`
	GoogleBloggerConfig   *GoogleBloggerConfig   `yaml:"google_blogger_config,omitempty" json:"google_blogger_config,omitempty"`
	ByteDanceDouYinConfig *ByteDanceDouYinConfig `yaml:"byte_dance_douyin_config,omitempty" json:"byte_dance_douyin_config,omitempty"`
	RedBookJuGuangConfig  *RedBookJuGuangConfig  `yaml:"redbook_juguang_config,omitempty" json:"redbook_juguang_config,omitempty"`
	BiliBiliConfig        *BiliBiliConfig        `yaml:"bilibili_config,omitempty" json:"bilibili_config,omitempty"`
	LegacyBiliBiliConfig  *BiliBiliConfig        `yaml:"bilbili_config,omitempty" json:"-"`
	CustomConfig          map[string]any         `yaml:"custom_config,omitempty" json:"custom_config,omitempty"`
	Meta                  map[string]string      `yaml:"meta,omitempty" json:"meta,omitempty"`
}

// FirstSelection 返回第一个可用的 Provider/App/Mode（用于默认值）
func (cfg *AccessTokenProvidersConfig) FirstSelection() (*AccessTokenProvider, *AccessTokenProviderApp, *AccessTokenAuthMode) {
	if cfg == nil {
		return nil, nil, nil
	}
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if mode := app.DefaultMode(); mode != nil {
				return provider, app, mode
			}
		}
	}
	return nil, nil, nil
}

// FindProvider 返回指定 code 的 Provider 分组
func (cfg *AccessTokenProvidersConfig) FindProvider(code string) *AccessTokenProvider {
	if cfg == nil {
		return nil
	}
	target := strings.TrimSpace(code)
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		if target == "" || strings.EqualFold(provider.Code, target) {
			return provider
		}
	}
	return nil
}

// FindApp 查找 Provider 分组与 App code 对应的配置
func (cfg *AccessTokenProvidersConfig) FindApp(groupCode, appCode string) (*AccessTokenProvider, *AccessTokenProviderApp) {
	if cfg == nil {
		return nil, nil
	}
	groupCode = strings.TrimSpace(groupCode)
	appCode = strings.TrimSpace(appCode)
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		if groupCode != "" && !strings.EqualFold(provider.Code, groupCode) {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if appCode == "" || strings.EqualFold(app.Code, appCode) {
				return provider, app
			}
		}
	}
	return nil, nil
}

// FindAppByProviderCode 使用 provider_code（如 google_youtube）查找 App
func (cfg *AccessTokenProvidersConfig) FindAppByProviderCode(providerCode string) (*AccessTokenProvider, *AccessTokenProviderApp) {
	if cfg == nil {
		return nil, nil
	}
	target := strings.TrimSpace(providerCode)
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if target == "" || strings.EqualFold(app.ProviderCodeValue(), target) {
				return provider, app
			}
		}
	}
	return nil, nil
}

// ProviderCodeValue 返回 App 对应的 provider_code（优先自身，再退回 group+code）
func (app *AccessTokenProviderApp) ProviderCodeValue() string {
	if app == nil {
		return ""
	}
	if code := strings.TrimSpace(app.ProviderCode); code != "" {
		return code
	}
	return strings.TrimSpace(app.Code)
}

// DefaultMode 返回 App 的默认授权模式
func (app *AccessTokenProviderApp) DefaultMode() *AccessTokenAuthMode {
	if app == nil || len(app.AuthModes) == 0 {
		return nil
	}
	return app.AuthModes[0]
}

// FindMode 根据 key 查找授权模式，找不到则返回默认模式
func (app *AccessTokenProviderApp) FindMode(key string) *AccessTokenAuthMode {
	if app == nil {
		return nil
	}
	target := strings.TrimSpace(key)
	if target == "" {
		return app.DefaultMode()
	}
	for _, mode := range app.AuthModes {
		if mode == nil {
			continue
		}
		if strings.EqualFold(mode.Key, target) {
			return mode
		}
	}
	return app.DefaultMode()
}

// EffectiveProviderCode 返回模式对应的 provider_code（优先 mode，再退回 app）
func (mode *AccessTokenAuthMode) EffectiveProviderCode(app *AccessTokenProviderApp) string {
	if mode == nil {
		if app == nil {
			return ""
		}
		return app.ProviderCodeValue()
	}
	if code := strings.TrimSpace(mode.ProviderCode); code != "" {
		return code
	}
	if app == nil {
		return ""
	}
	return app.ProviderCodeValue()
}

// ConfigKind 返回当前授权模式绑定的配置类型（例如 google_youtube）
func (mode *AccessTokenAuthMode) ConfigKind() string {
	if mode == nil {
		return ""
	}
	switch {
	case mode.GoogleYouTubeConfig != nil:
		return "google_youtube"
	case mode.GoogleBloggerConfig != nil:
		return "google_blogger"
	case mode.ByteDanceDouYinConfig != nil:
		return "byte_dance_douyin"
	case mode.RedBookJuGuangConfig != nil:
		return "redbook_juguang"
	case mode.BiliConfig() != nil:
		return "bilbili"
	default:
		return ""
	}
}

// ConfigObject 返回绑定的具体配置指针
func (mode *AccessTokenAuthMode) ConfigObject() any {
	if mode == nil {
		return nil
	}
	switch mode.ConfigKind() {
	case "google_youtube":
		return mode.GoogleYouTubeConfig
	case "google_blogger":
		return mode.GoogleBloggerConfig
	case "byte_dance_douyin":
		return mode.ByteDanceDouYinConfig
	case "redbook_juguang":
		return mode.RedBookJuGuangConfig
	case "bilbili":
		return mode.BiliConfig()
	default:
		return nil
	}
}

// ClientConfig 返回模式对应的 ClientConfig
func (mode *AccessTokenAuthMode) ClientConfig() *ClientConfig {
	if mode == nil {
		return nil
	}
	switch mode.ConfigKind() {
	case "google_youtube":
		if mode.GoogleYouTubeConfig != nil {
			return mode.GoogleYouTubeConfig.ClientConfig
		}
	case "google_blogger":
		if mode.GoogleBloggerConfig != nil {
			return mode.GoogleBloggerConfig.ClientConfig
		}
	case "byte_dance_douyin":
		if mode.ByteDanceDouYinConfig != nil {
			return mode.ByteDanceDouYinConfig.ClientConfig
		}
	case "redbook_juguang":
		if mode.RedBookJuGuangConfig != nil {
			return mode.RedBookJuGuangConfig.ClientConfig
		}
	case "bilbili":
		if cfg := mode.BiliConfig(); cfg != nil {
			return cfg.ClientConfig
		}
	}
	return nil
}

func (mode *AccessTokenAuthMode) BiliConfig() *BiliBiliConfig {
	if mode == nil {
		return nil
	}
	if mode.BiliBiliConfig != nil {
		return mode.BiliBiliConfig
	}
	return mode.LegacyBiliBiliConfig
}

// ClientTokenProvidersConfig 描述 ClientToken Provider（例如微信）配置
type ClientTokenProvidersConfig struct {
	Providers []*ClientTokenProvider `yaml:"providers" json:"providers"`
}

// ClientTokenProvider 代表 ClientToken Provider 分组
type ClientTokenProvider struct {
	Code string                    `yaml:"code" json:"code"`
	Name string                    `yaml:"name" json:"name"`
	Apps []*ClientTokenProviderApp `yaml:"apps" json:"apps"`
}

// ClientTokenProviderApp 表示 Provider 下的具体 App
type ClientTokenProviderApp struct {
	Code         string                 `yaml:"code" json:"code"`
	Name         string                 `yaml:"name" json:"name"`
	ProviderCode string                 `yaml:"provider_code,omitempty" json:"provider_code,omitempty"`
	AuthModes    []*ClientTokenAuthMode `yaml:"auth_modes" json:"auth_modes"`
}

// ClientTokenAuthMode 描述 ClientToken 授权模式
type ClientTokenAuthMode struct {
	Key                         string                     `yaml:"key" json:"key"`
	Label                       string                     `yaml:"label,omitempty" json:"label,omitempty"`
	ProviderCode                string                     `yaml:"provider_code,omitempty" json:"provider_code,omitempty"`
	WechatOfficialAccountConfig *ClientTokenProviderConfig `yaml:"wechat_official_account_config,omitempty" json:"wechat_official_account_config,omitempty"`
	CustomConfig                map[string]any             `yaml:"custom_config,omitempty" json:"custom_config,omitempty"`
	Meta                        map[string]string          `yaml:"meta,omitempty" json:"meta,omitempty"`
}

// FirstSelection 返回默认 Provider/App/Mode
func (cfg *ClientTokenProvidersConfig) FirstSelection() (*ClientTokenProvider, *ClientTokenProviderApp, *ClientTokenAuthMode) {
	if cfg == nil {
		return nil, nil, nil
	}
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if mode := app.DefaultMode(); mode != nil {
				return provider, app, mode
			}
		}
	}
	return nil, nil, nil
}

// FindApp 根据分组与 app code 查找
func (cfg *ClientTokenProvidersConfig) FindApp(groupCode, appCode string) (*ClientTokenProvider, *ClientTokenProviderApp) {
	if cfg == nil {
		return nil, nil
	}
	groupCode = strings.TrimSpace(groupCode)
	appCode = strings.TrimSpace(appCode)
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		if groupCode != "" && !strings.EqualFold(provider.Code, groupCode) {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if appCode == "" || strings.EqualFold(app.Code, appCode) {
				return provider, app
			}
		}
		if groupCode != "" {
			break
		}
	}
	return nil, nil
}

// FindAppByProviderCode 根据 provider_code 查找 app
func (cfg *ClientTokenProvidersConfig) FindAppByProviderCode(providerCode string) (*ClientTokenProvider, *ClientTokenProviderApp) {
	if cfg == nil {
		return nil, nil
	}
	target := strings.TrimSpace(providerCode)
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if target == "" || strings.EqualFold(app.ProviderCodeValue(), target) {
				return provider, app
			}
		}
	}
	return nil, nil
}

// ProviderCodeValue 返回 app 的 provider_code
func (app *ClientTokenProviderApp) ProviderCodeValue() string {
	if app == nil {
		return ""
	}
	if code := strings.TrimSpace(app.ProviderCode); code != "" {
		return code
	}
	return strings.TrimSpace(app.Code)
}

// DefaultMode 返回默认授权模式
func (app *ClientTokenProviderApp) DefaultMode() *ClientTokenAuthMode {
	if app == nil || len(app.AuthModes) == 0 {
		return nil
	}
	return app.AuthModes[0]
}

// FindMode 根据 key 查找模式
func (app *ClientTokenProviderApp) FindMode(key string) *ClientTokenAuthMode {
	if app == nil {
		return nil
	}
	target := strings.TrimSpace(key)
	if target == "" {
		return app.DefaultMode()
	}
	for _, mode := range app.AuthModes {
		if mode == nil {
			continue
		}
		if strings.EqualFold(mode.Key, target) {
			return mode
		}
	}
	return app.DefaultMode()
}

// SessionTokenProvidersConfig 描述 SessionToken Provider 配置
type SessionTokenProvidersConfig struct {
	Providers []*SessionTokenProvider `yaml:"providers" json:"providers"`
}

// SessionTokenProvider 代表 SessionToken Provider 分组
type SessionTokenProvider struct {
	Code string                     `yaml:"code" json:"code"`
	Name string                     `yaml:"name" json:"name"`
	Apps []*SessionTokenProviderApp `yaml:"apps" json:"apps"`
}

// SessionTokenProviderApp 表示 Provider 下的 App
type SessionTokenProviderApp struct {
	Code         string                  `yaml:"code" json:"code"`
	Name         string                  `yaml:"name" json:"name"`
	ProviderCode string                  `yaml:"provider_code,omitempty" json:"provider_code,omitempty"`
	AuthModes    []*SessionTokenAuthMode `yaml:"auth_modes" json:"auth_modes"`
}

// SessionTokenAuthMode 表示 SessionToken 授权模式
type SessionTokenAuthMode struct {
	Key                     string                   `yaml:"key" json:"key"`
	Label                   string                   `yaml:"label,omitempty" json:"label,omitempty"`
	ProviderCode            string                   `yaml:"provider_code,omitempty" json:"provider_code,omitempty"`
	ZhihuSessionTokenConfig *ZhihuSessionTokenConfig `yaml:"zhihu_session_token_config,omitempty" json:"zhihu_session_token_config,omitempty"`
	CustomConfig            map[string]any           `yaml:"custom_config,omitempty" json:"custom_config,omitempty"`
	Meta                    map[string]string        `yaml:"meta,omitempty" json:"meta,omitempty"`
}

// FirstSelection 返回第一个可用的 SessionToken Provider/App/Mode
func (cfg *SessionTokenProvidersConfig) FirstSelection() (*SessionTokenProvider, *SessionTokenProviderApp, *SessionTokenAuthMode) {
	if cfg == nil {
		return nil, nil, nil
	}
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if mode := app.DefaultMode(); mode != nil {
				return provider, app, mode
			}
		}
	}
	return nil, nil, nil
}

// FindApp 查找 Provider/App
func (cfg *SessionTokenProvidersConfig) FindApp(groupCode, appCode string) (*SessionTokenProvider, *SessionTokenProviderApp) {
	if cfg == nil {
		return nil, nil
	}
	groupCode = strings.TrimSpace(groupCode)
	appCode = strings.TrimSpace(appCode)
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		if groupCode != "" && !strings.EqualFold(provider.Code, groupCode) {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if appCode == "" || strings.EqualFold(app.Code, appCode) {
				return provider, app
			}
		}
		if groupCode != "" {
			break
		}
	}
	return nil, nil
}

// FindAppByProviderCode 使用 provider_code 查找 App
func (cfg *SessionTokenProvidersConfig) FindAppByProviderCode(providerCode string) (*SessionTokenProvider, *SessionTokenProviderApp) {
	if cfg == nil {
		return nil, nil
	}
	target := strings.TrimSpace(providerCode)
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			if target == "" || strings.EqualFold(app.ProviderCodeValue(), target) {
				return provider, app
			}
		}
	}
	return nil, nil
}

// ProviderCodeValue 返回 App 的 provider_code
func (app *SessionTokenProviderApp) ProviderCodeValue() string {
	if app == nil {
		return ""
	}
	if code := strings.TrimSpace(app.ProviderCode); code != "" {
		return code
	}
	return strings.TrimSpace(app.Code)
}

// DefaultMode 返回默认 SessionToken 授权模式
func (app *SessionTokenProviderApp) DefaultMode() *SessionTokenAuthMode {
	if app == nil || len(app.AuthModes) == 0 {
		return nil
	}
	return app.AuthModes[0]
}

// FindMode 根据 key 查找授权模式
func (app *SessionTokenProviderApp) FindMode(key string) *SessionTokenAuthMode {
	if app == nil {
		return nil
	}
	target := strings.TrimSpace(key)
	if target == "" {
		return app.DefaultMode()
	}
	for _, mode := range app.AuthModes {
		if mode == nil {
			continue
		}
		if strings.EqualFold(mode.Key, target) {
			return mode
		}
	}
	return app.DefaultMode()
}

// ResolveWechatOfficialAccount 查找微信公众号配置
func (cfg *ClientTokenProvidersConfig) ResolveWechatOfficialAccount(providerCode, appCode, modeKey string) (*ClientTokenProviderConfig, error) {
	if cfg == nil {
		return nil, fmt.Errorf("client_token_providers 未配置")
	}
	if mode := resolveClientTokenMode(cfg, providerCode, appCode, modeKey); mode != nil && mode.WechatOfficialAccountConfig != nil {
		return mode.WechatOfficialAccountConfig, nil
	}
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			for _, mode := range app.AuthModes {
				if mode != nil && mode.WechatOfficialAccountConfig != nil {
					return mode.WechatOfficialAccountConfig, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("client_token_providers: 未找到 wechat_official_account_config")
}

func resolveClientTokenMode(cfg *ClientTokenProvidersConfig, providerCode, appCode, modeKey string) *ClientTokenAuthMode {
	if cfg == nil {
		return nil
	}
	code := strings.TrimSpace(providerCode)
	appCode = strings.TrimSpace(appCode)
	modeKey = strings.TrimSpace(modeKey)
	var app *ClientTokenProviderApp
	if code != "" {
		_, app = cfg.FindAppByProviderCode(code)
	}
	if app == nil && appCode != "" {
		_, app = cfg.FindApp("", appCode)
	}
	if app == nil {
		_, _, mode := cfg.FirstSelection()
		return mode
	}
	return app.FindMode(modeKey)
}

// ResolveZhihuSessionToken 查找知乎 SessionToken 配置
func (cfg *SessionTokenProvidersConfig) ResolveZhihuSessionToken(providerCode, appCode, modeKey string) (*ZhihuSessionTokenConfig, error) {
	if cfg == nil {
		return nil, fmt.Errorf("session_token_providers 未配置")
	}
	if mode := resolveSessionTokenMode(cfg, providerCode, appCode, modeKey); mode != nil && mode.ZhihuSessionTokenConfig != nil {
		return mode.ZhihuSessionTokenConfig, nil
	}
	for _, provider := range cfg.Providers {
		if provider == nil {
			continue
		}
		for _, app := range provider.Apps {
			if app == nil {
				continue
			}
			for _, mode := range app.AuthModes {
				if mode != nil && mode.ZhihuSessionTokenConfig != nil {
					return mode.ZhihuSessionTokenConfig, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("session_token_providers: 未找到知乎 SessionToken 配置")
}

func resolveSessionTokenMode(cfg *SessionTokenProvidersConfig, providerCode, appCode, modeKey string) *SessionTokenAuthMode {
	if cfg == nil {
		return nil
	}
	code := strings.TrimSpace(providerCode)
	appCode = strings.TrimSpace(appCode)
	modeKey = strings.TrimSpace(modeKey)
	var app *SessionTokenProviderApp
	if code != "" {
		_, app = cfg.FindAppByProviderCode(code)
	}
	if app == nil && appCode != "" {
		_, app = cfg.FindApp("", appCode)
	}
	if app == nil {
		_, _, mode := cfg.FirstSelection()
		return mode
	}
	return app.FindMode(modeKey)
}
