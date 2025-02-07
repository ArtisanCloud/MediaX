package config

type WeChatConfig struct {
	AppConfig

	ComponentAppID    string `yaml:"component_app_id" json:"componentAppId"`
	ComponentAppToken string `yaml:"component_app_token" json:"componentAppToken"`
}
