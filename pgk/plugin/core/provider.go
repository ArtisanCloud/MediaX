// Path: plugins/core/provider.go
package core

// PublishRequest 定义了发布请求的基本内容
type PublishRequest struct {
	Title   string
	Content string
}

// PublishResult 定义了发布的结果结构体
type PublishResult struct {
	Status  string // 可以是 "success", "failed", "in-progress" 等
	Message string
}

type Provider interface {
	// Initialize 用于初始化 Provider 配置
	Initialize(config map[string]interface{}) error

	// Publish 用于发布内容，支持可变参数和回调函数
	Publish(req PublishRequest, extraParams ...interface{}) (PublishResult, error)

	// Name 返回 Provider 的名称
	Name() string
}
