package callback

import (
	"errors"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	callbackpkg "github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken/callback"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// NewDispatcher 基于知乎配置创建 HTTP 回调派发器。
func NewDispatcher(cfg *config.ZhihuSessionTokenCallbackConfig, log *logger.Logger) (*callbackpkg.HTTPDispatcher, error) {
	if cfg == nil {
		return nil, errors.New("zhihu.sessiontoken: callback config is nil")
	}
	secret := cfg.ResolveSecret()
	if secret == "" {
		return nil, errors.New("zhihu.sessiontoken: callback secret is empty")
	}
	var delays []time.Duration
	for _, second := range cfg.RetryBackoff {
		if second <= 0 {
			continue
		}
		delays = append(delays, time.Duration(second)*time.Second)
	}
	options := []callbackpkg.HTTPDispatcherOption{}
	if len(delays) > 0 {
		options = append(options, callbackpkg.WithRetrySchedule(delays))
	}
	return callbackpkg.NewHTTPDispatcher(secret, log, options...), nil
}
