// Package config 提供小红书聚光平台相关的配置定义
//
// 小红书聚光平台是面向品牌主、MCN机构的内容营销平台
// 支持内容发布、数据分析、账号管理等功能
//
// 主要功能：
// 1. 笔记管理：发布、编辑、删除笔记
// 2. 数据分析：获取笔记数据、账号数据等
// 3. 账号管理：授权管理、基础信息等
//
// 详细文档请参考：
// https://open.xiaohongshu.com/document
package config

import "github.com/ArtisanCloud/MediaXCore/utils/object"

// RedBookJuGuangConfig 小红书聚光平台客户端配置
// 继承自 ClientConfig，使用 OAuth2.0 认证方式
type RedBookJuGuangConfig struct {
	*ClientConfig `yaml:",inline"` // 基础客户端配置，包含 API 地址、超时设置等

	// GetOAuthToken 获取 OAuth Token 的回调函数
	// 参数：
	//   - key: OAuth 密钥
	//   - refresh: 是否刷新 token
	// 返回值：
	//   - token: 包含 access_token 等信息的 HashMap
	GetOAuthToken func(key string, refresh bool) (token object.HashMap) `yaml:"token;omitempty" json:"token;omitempty"`
}
