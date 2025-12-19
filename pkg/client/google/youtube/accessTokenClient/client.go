package accessTokenClient

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	v3 "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/v3"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// GoogleYouTubeACClient 当前默认采用 v3 实现，后续若有 v4/GraphQL 可在此做版本路由。
type GoogleYouTubeACClient = v3.GoogleYouTubeACClient

// NewGoogleYouTubeACClient 暂时只支持 v3。
func NewGoogleYouTubeACClient(cfg *config.GoogleYouTubeConfig, log *logger.Logger, cache cache.ICache) (*GoogleYouTubeACClient, error) {
	// 未来可通过 cfg.ApiVersion 切换不同实现。
	return v3.NewGoogleYouTubeACClient(cfg, log, cache)
}
