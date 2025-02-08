package client

import (
	"github.com/ArtisanCloud/MediaX/v1/internal/service/douyin"
	"github.com/ArtisanCloud/MediaX/v1/internal/service/redbook"
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat"
	"github.com/ArtisanCloud/MediaX/v1/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type MediaX struct {
	Logger *logger.Logger       // 全局 Logger
	Cache  cache.CacheInterface // 全局 Cache
}

// NewMediaX 初始化 MediaX，Logger 和 Cache 是全局共享的
func NewMediaX(config *config.MediaXConfig, cache cache.CacheInterface) *MediaX {
	l := logger.NewLogger(config.Logger)
	return &MediaX{
		Logger: l,
		Cache:  cache,
	}
}

// CreateWechat 创建 WechatClient，支持传入 WeChat 配置
func (m *MediaX) CreateWechat(cfg *config.WeChatConfig) (*wechat.WeChatService, error) {
	return wechat.NewWeChatService(cfg, m.Logger, m.Cache)
}

// CreateDouYin 创建 DouYinClient，支持传入 DouYin 配置
func (m *MediaX) CreateDouYin(cfg *config.DouYinConfig) (*douyin.DouYinService, error) {
	return douyin.NewDouYinService(cfg, m.Logger, m.Cache)
}

// CreateRedBook 创建 RedBookClient，支持传入 RedBook 配置
func (m *MediaX) CreateRedBook(cfg *config.RedBookConfig) (*redbook.RedBookService, error) {
	return redbook.NewRedBookService(cfg, m.Logger, m.Cache)
}
