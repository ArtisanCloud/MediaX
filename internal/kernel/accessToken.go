package kernel

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	request2 "github.com/ArtisanCloud/MediaX/internal/kernel/request"
	"github.com/ArtisanCloud/MediaX/internal/kernel/response"
	response2 "github.com/ArtisanCloud/MediaX/internal/kernel/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/contract"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/helper"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
	"net/http"
	"time"
)

type AccessTokenHandler struct {
	HttpHelper *helper.RequestHelper
	Cache      cache.CacheInterface
	Logger     *logger.Logger
	Config     *config.AppConfig

	RequestMethod      string
	EndpointToGetToken string
	QueryName          string
	Token              *object.HashMap
	TokenKey           string
	CacheTokenKey      string
	CachePrefix        string

	GetCredentials func() *object.StringMap
	GetEndpoint    func() (string, error)

	SetCustomToken func(token *response.AccessTokenRes) interface{}
	GetCustomToken func(key string, refresh bool) object.HashMap

	GetMiddlewareOfLog func(l *logger.Logger) contract.RequestMiddleware
}

func NewAccessTokenHandler(cfg *config.AppConfig, logger *logger.Logger, cache cache.CacheInterface) (*AccessTokenHandler, error) {
	h, err := helper.NewRequestHelper(&helper.Config{
		BaseUrl: cfg.BaseUri,
		ClientConfig: &contract.ClientConfig{
			Timeout:  time.Duration(cfg.Timeout * float64(time.Second)),
			ProxyURI: cfg.ProxyUri,
		},
	})
	if err != nil {
		return nil, err
	}

	handler := &AccessTokenHandler{
		HttpHelper: h,
		Logger:     logger,
		Cache:      cache,
		Config:     cfg,

		RequestMethod:      http.MethodGet,
		EndpointToGetToken: "",
		QueryName:          "",
		Token:              nil,
		TokenKey:           "access_token",
		CachePrefix:        "mediax.access_token.",
	}

	handler.OverrideMethods()
	handler.RegisterHttpMiddlewares()

	return handler, nil
}

func (acHandler *AccessTokenHandler) OverrideMethods() {
	acHandler.OverrideGetEndpoint()
	acHandler.OverrideGetMiddlewareOfLog()
}

func (acHandler *AccessTokenHandler) OverrideGetEndpoint() {
	acHandler.GetEndpoint = func() (string, error) {
		if acHandler.EndpointToGetToken == "" {
			return "", errors.New("no endpoint for access token request")
		}

		return acHandler.EndpointToGetToken, nil
	}
}

func (acHandler *AccessTokenHandler) OverrideGetMiddlewareOfLog() {
	acHandler.GetMiddlewareOfLog = func(l *logger.Logger) contract.RequestMiddleware {
		return func(handle contract.RequestHandle) contract.RequestHandle {
			return func(request *http.Request) (response *http.Response, err error) {
				l = l.WithContext(request.Context())

				request2.LogRequest(l, request)
				response, err = handle(request)
				if err == nil {
					l.WithContext(request.Context())
					response2.LogResponse(l, response)
				}
				return response, err
			}
		}
	}
}

func (acHandler *AccessTokenHandler) RegisterHttpMiddlewares() {
	// log
	logMiddleware := acHandler.GetMiddlewareOfLog

	acHandler.HttpHelper.WithMiddleware(
		logMiddleware(acHandler.Logger),
		helper.HttpDebugMiddleware(acHandler.Config.HttpDebug),
	)
}

func (acHandler *AccessTokenHandler) GetDefaultCacheKey() string {
	credentials := *acHandler.GetCredentials()
	data := fmt.Sprintf("%s%s%s", credentials["appid"], credentials["secret"], credentials["neededText"])
	buffer := md5.Sum([]byte(data))
	cacheKey := acHandler.CachePrefix + hex.EncodeToString(buffer[:])

	return cacheKey
}

func (acHandler *AccessTokenHandler) SetCacheKey(key string) {
	acHandler.CacheTokenKey = key
}

func (acHandler *AccessTokenHandler) GetCacheKey() string {
	cacheKey := ""
	if acHandler.CacheTokenKey != "" {
		cacheKey = acHandler.CacheTokenKey
	} else {
		cacheKey = acHandler.GetDefaultCacheKey()
	}

	return cacheKey
}

func (acHandler *AccessTokenHandler) getFormatToken(token object.HashMap) (*response.AccessTokenRes, error) {
	res := &response.AccessTokenRes{}
	err := object.HashMapToStructure(&token, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (acHandler *AccessTokenHandler) GetRefreshedToken() (*response2.AccessTokenRes, error) {
	return acHandler.GetToken(context.Background(), true)
}

func (acHandler *AccessTokenHandler) Refresh(ctx context.Context) *AccessTokenHandler {
	acHandler.GetToken(ctx, true)

	return acHandler
}

func (acHandler *AccessTokenHandler) sendRequest(ctx context.Context, credential *object.StringMap) (*response.AccessTokenRes, error) {
	key := "json"
	if acHandler.RequestMethod == http.MethodGet {
		key = "query"
	}
	options := &object.HashMap{
		key: credential,
	}

	res := &response.AccessTokenRes{}

	strEndpoint, err := acHandler.GetEndpoint()
	if err != nil {
		return nil, err
	}

	df := acHandler.HttpHelper.Df().WithContext(ctx).Uri(strEndpoint).
		Method(acHandler.RequestMethod)

	// 检查是否需要有请求参数配置
	// set query key values
	if (*options)["query"] != nil {
		queries := (*options)["query"].(*object.StringMap)
		if queries != nil {
			for k, v := range *queries {
				df.Query(k, v)
			}
		}
	}

	// set body json
	if (*options)["json"] != nil {
		df.Json((*options)["json"])
	}
	//if (*options)["form_params"] != nil {
	//	df.Json((*options)["form_params"])
	//}

	rs, err := df.Request()
	if err != nil {
		return nil, err
	}

	// decode response body to outBody
	err = acHandler.HttpHelper.ParseResponseBodyContent(rs, res)

	return res, err
}

func (acHandler *AccessTokenHandler) SetToken(ctx context.Context, token *response.AccessTokenRes) (acToken *AccessTokenHandler, err error) {
	if token.ExpiresIn <= 0 {
		token.ExpiresIn = 7200
	}

	// convert token to hashmap
	tokenJson, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	// set token into cache
	if acHandler.SetCustomToken != nil {
		customToken := acHandler.SetCustomToken(token)
		err = acHandler.Cache.Set(ctx, acHandler.GetCacheKey(), customToken, time.Duration(token.ExpiresIn)*time.Second)
	} else {
		err = acHandler.Cache.Set(ctx, acHandler.GetCacheKey(), tokenJson, time.Duration(token.ExpiresIn)*time.Second)
	}

	if err != nil {
		return nil, err
	}
	exist, err := acHandler.Cache.Exists(ctx, acHandler.GetCacheKey())
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("failed to cache access token")
	}
	return acHandler, err
}

func (acHandler *AccessTokenHandler) GetToken(ctx context.Context, refresh bool) (resToken *response.AccessTokenRes, err error) {
	cacheKey := acHandler.GetCacheKey()

	// 如果客户有中控的场景，可以由客户自己提供token的方法
	if acHandler.GetCustomToken != nil {
		token := acHandler.GetCustomToken(cacheKey, refresh)
		return acHandler.getFormatToken(token)
	}

	// get token from cache
	exist, err := acHandler.Cache.Exists(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if !refresh && exist {
		value, err := acHandler.Cache.Get(ctx, cacheKey)
		if err == nil && value != nil {
			resToken := &response.AccessTokenRes{}
			err = json.Unmarshal(value, resToken)
			return resToken, err
		}
	}

	// request token from power
	resToken, err = acHandler.sendRequest(ctx, acHandler.GetCredentials())
	if err != nil {
		return nil, err
	}
	if resToken.AccessToken == "" {
		return nil, fmt.Errorf("get access token error")
	}
	_, err = acHandler.SetToken(ctx, resToken)

	return resToken, err
}

func (acHandler *AccessTokenHandler) getQuery(ctx context.Context) (*object.StringMap, error) {
	// set the current token key
	var key string
	if acHandler.QueryName != "" {
		key = acHandler.QueryName
	} else {
		key = acHandler.TokenKey
	}

	// get token string power
	resToken, err := acHandler.GetToken(ctx, false)
	if err != nil {
		return nil, err
	}
	arrayReturn := &object.StringMap{
		key: resToken.AccessToken,
	}

	return arrayReturn, err
}

func (acHandler *AccessTokenHandler) ApplyToRequest(request *http.Request) (*http.Request, error) {
	// query Access Token power
	mapToken, err := acHandler.getQuery(request.Context())
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	for key, value := range *mapToken {
		q.Set(key, value)
	}
	request.URL.RawQuery = q.Encode()

	return request, err
}
