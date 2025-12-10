package config

// ZhihuConfig 聚合知乎相关配置。
type ZhihuConfig struct {
	SessionToken *ZhihuSessionTokenConfig `yaml:"sessionToken" json:"sessionToken"`
}

// ZhihuSessionTokenConfig 描述知乎 SessionToken 适配器所需的全部配置。
type ZhihuSessionTokenConfig struct {
	Service       ZhihuSessionTokenServiceConfig       `yaml:"service" json:"service"`
	Authenticator ZhihuSessionTokenAuthenticatorConfig `yaml:"authenticator" json:"authenticator"`
	Harvester     ZhihuSessionTokenHarvesterConfig     `yaml:"harvester" json:"harvester"`
	Callback      ZhihuSessionTokenCallbackConfig      `yaml:"callback" json:"callback"`
	Network       ZhihuSessionTokenNetworkConfig       `yaml:"network" json:"network"`
}

// ZhihuSessionTokenServiceConfig 描述服务端基础能力，例如 API Host、鉴权信息等。
type ZhihuSessionTokenServiceConfig struct {
	BaseURL   string `yaml:"base_url" json:"base_url"`
	APIToken  string `yaml:"api_token" json:"api_token"`
	Timeout   int    `yaml:"timeout" json:"timeout"` // 单位：秒
	HTTPDebug bool   `yaml:"http_debug,omitempty" json:"http_debug,omitempty"`
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
	CallbackSecret string `yaml:"callback_secret" json:"callback_secret"`
	MaxRetry       int    `yaml:"max_retry" json:"max_retry"`
	RetryBackoff   []int  `yaml:"retry_backoff" json:"retry_backoff"`
}

// ZhihuSessionTokenNetworkConfig 描述代理池、IP策略等网络限定。
type ZhihuSessionTokenNetworkConfig struct {
	ProxyPool      string `yaml:"proxy_pool" json:"proxy_pool"`
	IPStrategy     string `yaml:"ip_strategy" json:"ip_strategy"`
	RequestTimeout int    `yaml:"request_timeout" json:"request_timeout"` // 秒
}
