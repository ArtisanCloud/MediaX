package redbook

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/contract"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/helper"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"time"
)

type RedBookService struct {
	Logger     *logger.Logger        // 全局 Logger
	Cache      cache.ICache          // 全局 Cache
	HttpHelper *helper.RequestHelper // 全局 HttpClient
}

func NewRedBookService(cfg *config.RedBookConfig, logger *logger.Logger, cache cache.ICache) (*RedBookService, error) {

	httpHelper, err := helper.NewRequestHelper(&helper.Config{
		BaseUrl: cfg.ApiUrl,
		ClientConfig: &contract.ClientConfig{
			Timeout:  time.Duration(cfg.Timeout * float64(time.Second)),
			ProxyURI: cfg.ProxyApiUrl,
		},
	})

	if err != nil {
		return nil, err
	}

	return &RedBookService{
		Logger:     logger,
		Cache:      cache,
		HttpHelper: httpHelper,
	}, nil
}
