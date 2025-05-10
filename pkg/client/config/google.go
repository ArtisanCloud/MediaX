// Package config 提供 Google 平台相关的配置定义
package config

import "github.com/ArtisanCloud/MediaXCore/utils/object"

// GoogleYouTubeConfig YouTube 客户端配置
// 继承自 ClientConfig，添加了 YouTube 特有的配置项
type GoogleYouTubeConfig struct {
	*ClientConfig `yaml:",inline"` // 基础客户端配置

	// GetOAuthToken 获取 OAuth Token 的回调函数
	// 参数：
	//   - key: OAuth 密钥
	//   - refresh: 是否刷新 token
	// 返回值：
	//   - token: 包含 access_token 等信息的 HashMap
	GetOAuthToken func(key string, refresh bool) (token object.HashMap) `yaml:"token;omitempty" json:"token;omitempty"`

	// OauthKey OAuth 密钥
	// 用于标识和获取特定的 OAuth 配置
	OauthKey string `yaml:"oauth_key;omitempty" json:"oauth_key;omitempty"`
}

// GoogleBloggerConfig Blogger 客户端配置
// 继承自 ClientConfig，添加了 Blogger 特有的配置项
type GoogleBloggerConfig struct {
	*ClientConfig `yaml:",inline"` // 基础客户端配置

	// GetOAuthToken 获取 OAuth Token 的回调函数
	// 参数：
	//   - key: OAuth 密钥
	//   - refresh: 是否刷新 token
	// 返回值：
	//   - token: 包含 access_token 等信息的 HashMap
	GetOAuthToken func(key string, refresh bool) (token object.HashMap) `yaml:"token;omitempty" json:"token;omitempty"`

	// OauthKey OAuth 密钥
	// 用于标识和获取特定的 OAuth 配置
	OauthKey string `yaml:"oauth_key;omitempty" json:"oauth_key;omitempty"`
}
