// Package config 提供抖音开放平台相关的配置定义
//
// 抖音开放平台提供了两种认证方式：
// 1. ClientToken：适用于服务端API调用
// 2. AccessToken：适用于用户授权的API调用
//
// 详细文档请参考：
// https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/
package config

// ByteDanceDouYinConfig 抖音开放平台客户端配置
// 继承自 ClientConfig，支持 ClientToken 和 AccessToken 两种认证方式
type ByteDanceDouYinConfig struct {
	*ClientConfig `yaml:",inline"` // 基础客户端配置，包含 API 地址、超时设置等

	// GetOAuthToken 获取 OAuth Token 的回调函数（暂未启用）
	// 后续可能会添加用于处理用户授权 token 的相关功能
	// GetOAuthToken func(key string, refresh bool) (token object.HashMap) `yaml:"token;omitempty" json:"token;omitempty"`
}
