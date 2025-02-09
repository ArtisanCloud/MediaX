package config

type WeChatOfficialAccountConfig struct {
	AppConfig `yaml:"app" json:"app"`

	ComponentAppID    string `yaml:"component_app_id" json:"componentAppId"`
	ComponentAppToken string `yaml:"component_app_token" json:"componentAppToken"`
}
