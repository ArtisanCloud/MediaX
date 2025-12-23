// Package config 提供抖音开放平台相关的配置定义
//
// 抖音开放平台提供了两种认证方式：
// 1. ClientToken：适用于服务端API调用
// 2. AccessToken：适用于用户授权的API调用
//
// 详细文档请参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/
package config

import (
	"strings"

	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// ByteDanceDouYinConfig 抖音开放平台客户端配置
// 继承自 ClientConfig：`oauth` 段需要通过 `DOUYIN_CLIENT_ID`/`DOUYIN_CLIENT_SECRET`/`DOUYIN_SCOPE`/`DOUYIN_REDIRECT_URL`/`DOUYIN_OAUTH_URL`/`DOUYIN_ACCESS_TOKEN_URL` 等环境变量注入，保证仓库不落地明文。
// DouYin AccessToken Flow 仅支持单 app（`provider_app=douyin`，对应 `<app>=default`），Flow TTL 必须等于 DouYin OAuth 返回的 `expires_in`，刷新成功后同样需要更新 TTL。
type ByteDanceDouYinConfig struct {
	*ClientConfig `yaml:",inline"` // 基础客户端配置（API URL、代理、timeout、oauth），推荐只在 YAML 中引用环境变量或 `config.yaml` 的占位符

	// OauthKey 区分租户/账号的 key，将展示在 `/debug` Provider 卡片与 Flow 记录中，例如 `${DOUYIN_OAUTH_KEY:-default}`。
	OauthKey string `yaml:"oauth_key,omitempty" json:"oauth_key,omitempty"`

	// GetOAuthToken 允许直接注入 AccessToken/ExpiresIn（例如本地脚本调试），同 Google/Bili 结构保持一致。
	// 建议仅在测试环境通过 `DOUYIN_ACCESS_TOKEN` 这类环境变量填充，生产环境仍需走 OAuth Flow。
	GetOAuthToken func(key string, refresh bool) (token object.HashMap) `yaml:"token;omitempty" json:"token;omitempty"`

	// ClientToken、Cache、Redis、APIToken 用于 ClientToken 调试台/CLI 与 SDK 的共享配置。
	ClientToken *ByteDanceDouYinClientTokenCredential `yaml:"client_token,omitempty" json:"client_token,omitempty"`
	Cache       *ClientTokenCacheConfig               `yaml:"cache,omitempty" json:"cache,omitempty"`
	Redis       *ClientTokenRedisConfig               `yaml:"redis,omitempty" json:"redis,omitempty"`
	APIToken    string                                `yaml:"api_token,omitempty" json:"api_token,omitempty"`

	// 风控字段，仅当需要调用对 device_id/risk_info 有要求的 API 时启用；默认留空即禁用。
	DeviceID string `yaml:"device_id,omitempty" json:"device_id,omitempty"`
	RiskInfo string `yaml:"risk_info,omitempty" json:"risk_info,omitempty"`
}

// ByteDanceDouYinClientTokenCredential 记录 DouYin client_token 所需 client_key/client_secret
type ByteDanceDouYinClientTokenCredential struct {
	ClientKey    string `yaml:"client_key" json:"client_key"`
	ClientSecret string `yaml:"client_secret" json:"client_secret"`
}

// ClientTokenCredential 返回凭证结构
func (cfg *ByteDanceDouYinConfig) ClientTokenCredential() *ByteDanceDouYinClientTokenCredential {
	if cfg == nil {
		return nil
	}
	return cfg.ClientToken
}

// ClientKeyValue 返回 client_key
func (cred *ByteDanceDouYinClientTokenCredential) ClientKeyValue() string {
	if cred == nil {
		return ""
	}
	return strings.TrimSpace(cred.ClientKey)
}

// ClientSecretValue 返回 client_secret
func (cred *ByteDanceDouYinClientTokenCredential) ClientSecretValue() string {
	if cred == nil {
		return ""
	}
	return strings.TrimSpace(cred.ClientSecret)
}

// EnsureClientConfig 初始化 ClientConfig/Base/OAuth
func (cfg *ByteDanceDouYinConfig) EnsureClientConfig() *ClientConfig {
	if cfg == nil {
		return nil
	}
	if cfg.ClientConfig == nil {
		cfg.ClientConfig = &ClientConfig{}
	}
	if cfg.ClientConfig.BaseConfig == nil {
		cfg.ClientConfig.BaseConfig = &BaseConfig{}
	}
	if cfg.ClientConfig.OAuthConfig == nil {
		cfg.ClientConfig.OAuthConfig = &OAuthConfig{}
	}
	return cfg.ClientConfig
}

// NormalizeClientTokenCredentials 将 client_token 凭证写入 OAuth 配置，便于 SDK/调试台复用
func (cfg *ByteDanceDouYinConfig) NormalizeClientTokenCredentials() {
	cc := cfg.EnsureClientConfig()
	if cc == nil || cc.OAuthConfig == nil {
		return
	}
	cred := cfg.ClientTokenCredential()
	if cred == nil {
		return
	}
	if v := cred.ClientKeyValue(); v != "" {
		cc.OAuthConfig.ClientID = v
	}
	if v := cred.ClientSecretValue(); v != "" {
		cc.OAuthConfig.ClientSecret = v
	}
}

// EffectiveRedisKey 返回缓存 key，默认 clientToken:douyin:<client_key>
func (cfg *ByteDanceDouYinConfig) EffectiveRedisKey(defaultKey string) string {
	if cfg == nil {
		return ""
	}
	if cfg.Cache != nil {
		if key := strings.TrimSpace(cfg.Cache.RedisKey); key != "" {
			return key
		}
	}
	key := ""
	if cred := cfg.ClientTokenCredential(); cred != nil {
		key = cred.ClientKeyValue()
	}
	if key == "" {
		key = strings.TrimSpace(defaultKey)
	}
	if key == "" {
		return ""
	}
	return "clientToken:douyin:" + key
}

// EffectiveTTLSeconds 返回缓存 TTL，默认 7000 秒
func (cfg *ByteDanceDouYinConfig) EffectiveTTLSeconds() int {
	if cfg == nil || cfg.Cache == nil || cfg.Cache.TTLSeconds <= 0 {
		return 7000
	}
	return cfg.Cache.TTLSeconds
}

// EffectiveRefreshBefore 返回自动刷新阈值，默认 600 秒
func (cfg *ByteDanceDouYinConfig) EffectiveRefreshBefore() int {
	if cfg == nil || cfg.Cache == nil || cfg.Cache.RefreshBeforeSeconds <= 0 {
		return 600
	}
	return cfg.Cache.RefreshBeforeSeconds
}

func (cfg *ByteDanceDouYinConfig) hasDeviceID() bool {
	return cfg != nil && strings.TrimSpace(cfg.DeviceID) != ""
}

func (cfg *ByteDanceDouYinConfig) hasRiskInfo() bool {
	return cfg != nil && strings.TrimSpace(cfg.RiskInfo) != ""
}

// RiskControlActivated 当用户尝试配置 device_id/risk_info 任意字段时返回 true。
func (cfg *ByteDanceDouYinConfig) RiskControlActivated() bool {
	return cfg != nil && (cfg.hasDeviceID() || cfg.hasRiskInfo())
}

// RiskControlReady 当 device_id 与 risk_info 均配置完成时返回 true。
func (cfg *ByteDanceDouYinConfig) RiskControlReady() bool {
	return cfg != nil && cfg.hasDeviceID() && cfg.hasRiskInfo()
}

// MissingRiskControlFields 返回缺失的风控字段列表。
func (cfg *ByteDanceDouYinConfig) MissingRiskControlFields() []string {
	if cfg == nil {
		return []string{"device_id", "risk_info"}
	}
	missing := []string{}
	if !cfg.hasDeviceID() {
		missing = append(missing, "device_id")
	}
	if !cfg.hasRiskInfo() {
		missing = append(missing, "risk_info")
	}
	return missing
}
