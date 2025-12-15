package config

import "strings"

// ZhihuConfig 聚合知乎相关配置。
type ZhihuConfig struct {
	SessionToken *ZhihuSessionTokenConfig `yaml:"sessiontoken" json:"sessiontoken"`
}

// ZhihuSessionTokenConfig 描述知乎 SessionToken 适配器所需的全部配置。
type ZhihuSessionTokenConfig struct {
	Strategy      string                               `yaml:"strategy,omitempty" json:"strategy,omitempty"` // 策略版本标识，如 zhihu_pc_v1
	Service       ZhihuSessionTokenServiceConfig       `yaml:"service" json:"service"`
	Authenticator ZhihuSessionTokenAuthenticatorConfig `yaml:"authenticator" json:"authenticator"`
	Harvester     ZhihuSessionTokenHarvesterConfig     `yaml:"harvester" json:"harvester"`
	Callback      ZhihuSessionTokenCallbackConfig      `yaml:"callback" json:"callback"`
	Network       ZhihuSessionTokenNetworkConfig       `yaml:"network" json:"network"`
	API           ZhihuSessionTokenAPIConfig           `yaml:"api" json:"api"`
}

// ZhihuSessionTokenServiceConfig 描述服务端基础能力，例如 API Host、鉴权信息等。
type ZhihuSessionTokenServiceConfig struct {
	BaseURL    string `yaml:"base_url" json:"base_url"`
	APIToken   string `yaml:"api_token" json:"api_token"`
	Timeout    int    `yaml:"timeout" json:"timeout"` // 单位：秒
	HTTPDebug  bool   `yaml:"http_debug,omitempty" json:"http_debug,omitempty"`
	APIVersion string `yaml:"api_version,omitempty" json:"api_version,omitempty"`
}

// ZhihuSessionTokenAuthenticatorConfig 配置不同入口、UA 与脚本策略。
type ZhihuSessionTokenAuthenticatorConfig struct {
	Entries          []ZhihuSessionTokenEntryConfig `yaml:"entries" json:"entries"`
	DefaultUserAgent string                         `yaml:"default_user_agent" json:"default_user_agent"`
	ScriptIDs        []string                       `yaml:"script_ids" json:"script_ids"`
	CaptchaStrategy  string                         `yaml:"captcha_strategy" json:"captcha_strategy"`
}

// ZhihuSessionTokenEntryConfig 描述单个登录入口。
type ZhihuSessionTokenEntryConfig struct {
	Type string `yaml:"type" json:"type"`
	URL  string `yaml:"url" json:"url"`
}

// ZhihuSessionTokenHarvesterConfig 描述需要监听的 cookie/header 以及脚本信息。
type ZhihuSessionTokenHarvesterConfig struct {
	WatchCookies    []string `yaml:"watch_cookies" json:"watch_cookies"`
	WatchHeaders    []string `yaml:"watch_headers" json:"watch_headers"`
	HarvestScriptID string   `yaml:"harvest_script_id" json:"harvest_script_id"`
}

// ZhihuSessionTokenCallbackConfig 配置回调签名、重试策略。
type ZhihuSessionTokenCallbackConfig struct {
	Secret         string `yaml:"secret,omitempty" json:"secret,omitempty"`
	CallbackSecret string `yaml:"callback_secret,omitempty" json:"callback_secret,omitempty"`
	MaxRetry       int    `yaml:"max_retry" json:"max_retry"`
	RetryBackoff   []int  `yaml:"retry_backoff" json:"retry_backoff"`
}

// ResolveSecret 返回有效的回调 secret。
func (c *ZhihuSessionTokenCallbackConfig) ResolveSecret() string {
	if c == nil {
		return ""
	}
	if secret := strings.TrimSpace(c.Secret); secret != "" {
		return secret
	}
	return strings.TrimSpace(c.CallbackSecret)
}

// ZhihuSessionTokenNetworkConfig 描述代理池、IP策略等网络限定。
type ZhihuSessionTokenNetworkConfig struct {
	Proxy          string `yaml:"proxy,omitempty" json:"proxy,omitempty"`
	ProxyPool      string `yaml:"proxy_pool,omitempty" json:"proxy_pool,omitempty"`
	IPStrategy     string `yaml:"ip_strategy" json:"ip_strategy"`
	RequestTimeout int    `yaml:"request_timeout" json:"request_timeout"` // 秒
}

// ZhihuSessionTokenAPIConfig 描述 API 层的可调参数。
type ZhihuSessionTokenAPIConfig struct {
	RetryBackoff []int `yaml:"retry_backoff" json:"retry_backoff"`
}
