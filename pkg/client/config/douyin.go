// Package config 提供抖音开放平台相关的配置定义
//
// 抖音开放平台提供了两种认证方式：
// 1. ClientToken：适用于服务端API调用
// 2. AccessToken：适用于用户授权的API调用
//
// 详细文档请参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/
package config

import "github.com/ArtisanCloud/MediaXCore/utils/object"

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
}
