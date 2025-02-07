package kernel

import (
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/contract"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/helper"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"time"
)

type BaseClient struct {
	HttpHelper *helper.RequestHelper
	Logger     *logger.Logger
	Cache      cache.CacheInterface

	BaseURI  string
	QueryRaw bool

	Token *AccessToken

	GetMiddlewareOfAccessToken        contract.RequestMiddleware
	GetMiddlewareOfLog                func(l *logger.Logger) contract.RequestMiddleware
	GetMiddlewareOfRefreshAccessToken func(retry int) contract.RequestMiddleware
}

func NewBaseClient(logger *logger.Logger, cache cache.CacheInterface (*BaseClient, error) {
	config := (*app).GetConfig()
	baseURI := config.GetString("http.base_uri", "/")
	proxyURI := config.GetString("http.proxy_uri", "")
	timeout := config.GetFloat64("http.timeout", 5)

	if token == nil {
		token = (*app).GetAccessToken()
	}

	h, err := helper.NewRequestHelper(&helper.Config{
		BaseUrl: baseURI,
		ClientConfig: &contract.ClientConfig{
			Timeout:  time.Duration(timeout * float64(time.Second)),
			ProxyURI: proxyURI,
		},
	})
	if err != nil {
		return nil, err
	}
	if proxyURL := config.GetString("http.proxy", ""); proxyURL != "" {
		h.GetClient()
	}
	client := &BaseClient{
		HttpHelper: h,
		App:        app,
		Token:      token,
	}
	client.Logger = (*client.App).GetComponent("Logger").(contract2.LoggerInterface)

	// to be setup middleware here
	client.OverrideGetMiddlewares()
	client.RegisterHttpMiddlewares()

	mchID := config.GetString("mch_id", "")
	serialNO := config.GetString("serial_no", "")
	keyPath := config.GetString("key_path", "")
	if mchID != "" && serialNO != "" && keyPath != "" {
		client.Signer = &support.SHA256WithRSASigner{
			MchID:               mchID,
			CertificateSerialNo: serialNO,
			PrivateKeyPath:      keyPath,
		}
	}

	return client, nil
}
