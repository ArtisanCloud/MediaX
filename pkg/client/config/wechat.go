package config

import "strings"

type WeChatOfficialAccountConfig struct {
	*ClientConfig `yaml:",inline"`

	//ComponentAppID    string `yaml:"component_app_id" json:"componentAppId"`
	//ComponentAppToken string `yaml:"component_app_token" json:"componentAppToken"`
}

// WechatClientTokenCredential 描述微信公众号 ClientToken 所需凭证
type WechatClientTokenCredential struct {
	AppID         string `yaml:"appid" json:"appid"`
	AppSecret     string `yaml:"appsecret" json:"appsecret"`
	MessageToken  string `yaml:"message_token,omitempty" json:"message_token,omitempty"`
	MessageAESKey string `yaml:"message_aes_key,omitempty" json:"message_aes_key,omitempty"`
}

// ClientTokenProviderConfig 代表调试 server 所需的完整配置，包含客户端配置与缓存策略
type ClientTokenProviderConfig struct {
	*WeChatOfficialAccountConfig `yaml:",inline"`
	ClientToken                  *WechatClientTokenCredential `yaml:"client_token,omitempty" json:"client_token,omitempty"`
	APIToken                     string                       `yaml:"api_token,omitempty" json:"api_token,omitempty"`
	Cache                        *ClientTokenCacheConfig      `yaml:"cache,omitempty" json:"cache,omitempty"`
	Redis                        *ClientTokenRedisConfig      `yaml:"redis,omitempty" json:"redis,omitempty"`
}

// ClientTokenCacheConfig 控制调试 server 存储 token 的 key/TTL
type ClientTokenCacheConfig struct {
	RedisKey             string `yaml:"redis_key,omitempty" json:"redis_key,omitempty"`
	TTLSeconds           int    `yaml:"ttl_seconds,omitempty" json:"ttl_seconds,omitempty"`
	RefreshBeforeSeconds int    `yaml:"refresh_before_seconds,omitempty" json:"refresh_before_seconds,omitempty"`
}

// ClientTokenRedisConfig 记录 Redis 连接信息
type ClientTokenRedisConfig struct {
	Addr     string `yaml:"addr,omitempty" json:"addr,omitempty"`
	DB       int    `yaml:"db,omitempty" json:"db,omitempty"`
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

// AppIDValue 返回归一化后的 AppID
func (cfg *WechatClientTokenCredential) AppIDValue() string {
	if cfg == nil {
		return ""
	}
	return strings.TrimSpace(cfg.AppID)
}

// AppSecretValue 返回 AppSecret
func (cfg *WechatClientTokenCredential) AppSecretValue() string {
	if cfg == nil {
		return ""
	}
	return strings.TrimSpace(cfg.AppSecret)
}

// MessageTokenValue 返回消息校验 token
func (cfg *WechatClientTokenCredential) MessageTokenValue() string {
	if cfg == nil {
		return ""
	}
	return strings.TrimSpace(cfg.MessageToken)
}

// MessageAESKeyValue 返回 AES Key
func (cfg *WechatClientTokenCredential) MessageAESKeyValue() string {
	if cfg == nil {
		return ""
	}
	return strings.TrimSpace(cfg.MessageAESKey)
}

// Credentials 返回凭证指针
func (cfg *ClientTokenProviderConfig) Credentials() *WechatClientTokenCredential {
	if cfg == nil {
		return nil
	}
	return cfg.ClientToken
}

// OfficialAccountConfig 返回嵌入的公众号客户端配置
func (cfg *ClientTokenProviderConfig) OfficialAccountConfig() *WeChatOfficialAccountConfig {
	if cfg == nil {
		return nil
	}
	if cfg.WeChatOfficialAccountConfig == nil {
		cfg.WeChatOfficialAccountConfig = &WeChatOfficialAccountConfig{}
	}
	return cfg.WeChatOfficialAccountConfig
}

// AppIDValue 返回凭证 AppID，若为空则退回 OAuth client_id
func (cfg *ClientTokenProviderConfig) AppIDValue() string {
	if cfg == nil {
		return ""
	}
	if cred := cfg.Credentials(); cred != nil && cred.AppIDValue() != "" {
		return cred.AppIDValue()
	}
	if cfg.WeChatOfficialAccountConfig != nil &&
		cfg.WeChatOfficialAccountConfig.ClientConfig != nil &&
		cfg.WeChatOfficialAccountConfig.ClientConfig.OAuthConfig != nil {
		return strings.TrimSpace(cfg.WeChatOfficialAccountConfig.ClientConfig.OAuthConfig.ClientID)
	}
	return ""
}

// EffectiveRedisKey 返回缓存 key，优先显式配置
func (cfg *ClientTokenProviderConfig) EffectiveRedisKey(defaultAppID string) string {
	if cfg == nil {
		return ""
	}
	if cfg.Cache != nil {
		if key := strings.TrimSpace(cfg.Cache.RedisKey); key != "" {
			return key
		}
	}
	appID := cfg.AppIDValue()
	if appID == "" {
		appID = strings.TrimSpace(defaultAppID)
	}
	if appID == "" {
		return ""
	}
	return "clientToken:wechat:" + appID
}

// EffectiveTTLSeconds 返回 TTL，默认为 7000 秒
func (cfg *ClientTokenProviderConfig) EffectiveTTLSeconds() int {
	if cfg == nil || cfg.Cache == nil || cfg.Cache.TTLSeconds <= 0 {
		return 7000
	}
	return cfg.Cache.TTLSeconds
}

// EffectiveRefreshBefore 返回刷新阈值，默认 600 秒
func (cfg *ClientTokenProviderConfig) EffectiveRefreshBefore() int {
	if cfg == nil || cfg.Cache == nil || cfg.Cache.RefreshBeforeSeconds <= 0 {
		return 600
	}
	return cfg.Cache.RefreshBeforeSeconds
}
