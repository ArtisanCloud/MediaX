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
	"strings"
	"time"
)

type TokenHandler struct {
	HttpHelper *helper.RequestHelper
	Cache      cache.ICache
	Logger     *logger.Logger
	Config     *config.ClientConfig

	RequestMethod      string
	EndpointToGetToken string
	QueryName          string
	Token              *object.HashMap
	TokenKey           string
	CacheTokenKey      string
	CachePrefix        string

	GetCredentials func() *object.StringMap
	GetEndpoint    func() (string, error)

	SetCustomToken func(token interface{}) interface{}
	GetCustomToken func(key string, refresh bool) object.HashMap
	GetTokenQuery  func(ctx context.Context) (*object.StringMap, error)

	GetMiddlewareOfLog func(l *logger.Logger) contract.RequestMiddleware
}

func NewTokenHandler(cfg *config.ClientConfig, logger *logger.Logger, cache cache.ICache) (*TokenHandler, error) {
	h, err := helper.NewRequestHelper(&helper.Config{
		BaseUrl: cfg.ApiUrl,
		ClientConfig: &contract.ClientConfig{
			Timeout:  time.Duration(cfg.Timeout * float64(time.Second)),
			ProxyURI: cfg.ProxyOAuthUrl,
		},
	})
	if err != nil {
		return nil, err
	}

	handler := &TokenHandler{
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

func (tHandler *TokenHandler) OverrideMethods() {
	tHandler.OverrideGetTokenQuery()
	tHandler.OverrideGetEndpoint()
	tHandler.OverrideGetMiddlewareOfLog()
}

func (tHandler *TokenHandler) OverrideGetEndpoint() {
	tHandler.GetEndpoint = func() (string, error) {
		if tHandler.EndpointToGetToken == "" {
			return "", errors.New("no endpoint for access token request")
		}

		return tHandler.EndpointToGetToken, nil
	}
}

func (tHandler *TokenHandler) OverrideGetMiddlewareOfLog() {
	tHandler.GetMiddlewareOfLog = func(l *logger.Logger) contract.RequestMiddleware {
		return func(handle contract.RequestHandle) contract.RequestHandle {
			return func(request *http.Request) (response *http.Response, err error) {
				l = l.WithContext(request.Context())

				request2.LogRequest(tHandler.Config.HttpDebug, l, request)
				response, err = handle(request)
				if err == nil {
					l.WithContext(request.Context())
					response2.LogResponse(tHandler.Config.HttpDebug, l, response)
				}
				return response, err
			}
		}
	}
}

func (tHandler *TokenHandler) RegisterHttpMiddlewares() {
	// log
	logMiddleware := tHandler.GetMiddlewareOfLog

	tHandler.HttpHelper.WithMiddleware(
		logMiddleware(tHandler.Logger),
		helper.HttpDebugMiddleware(tHandler.Config.HttpDebug),
	)
}

func (tHandler *TokenHandler) GetDefaultCacheKey() string {
	credentials := *tHandler.GetCredentials()
	var builder strings.Builder

	// 遍历 credentials map，拼接所有字段值
	for _, value := range credentials {
		builder.WriteString(value)
	}

	// 计算 MD5
	buffer := md5.Sum([]byte(builder.String()))
	cacheKey := tHandler.CachePrefix + hex.EncodeToString(buffer[:])

	return cacheKey
}

func (tHandler *TokenHandler) SetCacheKey(key string) {
	tHandler.CacheTokenKey = key
}

func (tHandler *TokenHandler) GetCacheKey() string {
	cacheKey := ""
	if tHandler.CacheTokenKey != "" {
		cacheKey = tHandler.CacheTokenKey
	} else {
		cacheKey = tHandler.GetDefaultCacheKey()
	}

	return cacheKey
}

func (tHandler *TokenHandler) getFormatToken(token object.HashMap) (*response.AccessTokenRes, error) {
	res := &response.AccessTokenRes{}
	err := object.HashMapToStructure(&token, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (tHandler *TokenHandler) GetRefreshedToken() (*response2.AccessTokenRes, error) {
	resToken := &response2.AccessTokenRes{}
	err := tHandler.GetToken(context.Background(), true, resToken)
	return resToken, err
}

func (tHandler *TokenHandler) Refresh(ctx context.Context) *TokenHandler {
	resToken := &response2.AccessTokenRes{}
	tHandler.GetToken(ctx, true, resToken)

	return tHandler
}

func (tHandler *TokenHandler) sendRequest(ctx context.Context, credential *object.StringMap) (*response.AccessTokenRes, error) {
	key := "json"
	if tHandler.RequestMethod == http.MethodGet {
		key = "query"
	}
	options := &object.HashMap{
		key: credential,
	}

	res := &response.AccessTokenRes{}

	strEndpoint, err := tHandler.GetEndpoint()
	if err != nil {
		return nil, err
	}

	df := tHandler.HttpHelper.Df().WithContext(ctx).Uri(strEndpoint).
		Method(tHandler.RequestMethod)

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
	err = tHandler.HttpHelper.ParseResponseBodyContent(rs, res)

	return res, err
}

func (tHandler *TokenHandler) SetToken(ctx context.Context, token interface{}, expiresIn float64) (acToken *TokenHandler, err error) {
	if expiresIn <= 0 {
		expiresIn = 7200
	}
	// convert token to hashmap
	tokenJson, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	// set token into cache
	if tHandler.SetCustomToken != nil {
		customToken := tHandler.SetCustomToken(token)
		err = tHandler.Cache.Set(ctx, tHandler.GetCacheKey(), customToken, time.Duration(expiresIn)*time.Second)
	} else {
		err = tHandler.Cache.Set(ctx, tHandler.GetCacheKey(), tokenJson, time.Duration(expiresIn)*time.Second)
	}

	if err != nil {
		return nil, err
	}
	exist, err := tHandler.Cache.Exists(ctx, tHandler.GetCacheKey())
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("failed to cache access token")
	}
	return tHandler, err
}

func (tHandler *TokenHandler) GetToken(ctx context.Context, refresh bool, resToken interface{}) (err error) {
	cacheKey := tHandler.GetCacheKey()

	// 如果客户有中控的场景，可以由客户自己提供token的方法
	if tHandler.GetCustomToken != nil {
		resToken = tHandler.GetCustomToken(cacheKey, refresh)
		return nil
	}

	// get token from cache
	exist, err := tHandler.Cache.Exists(ctx, cacheKey)
	if err != nil {
		return err
	}
	if !refresh && exist {
		value, err := tHandler.Cache.Get(ctx, cacheKey)
		if err == nil && value != nil {
			return nil
		}
	}

	// request token from provider auth token api
	newToken, err := tHandler.sendRequest(ctx, tHandler.GetCredentials())
	if err != nil {
		return err
	}
	_, err = tHandler.SetToken(ctx, resToken, newToken.ExpiresIn)

	return err
}

func (tHandler *TokenHandler) OverrideGetTokenQuery() {
	tHandler.GetTokenQuery = func(ctx context.Context) (*object.StringMap, error) {
		// set the current token key
		var key string
		if tHandler.QueryName != "" {
			key = tHandler.QueryName
		} else {
			key = tHandler.TokenKey
		}

		// get token string power
		resToken := &response.AccessTokenRes{}
		err := tHandler.GetToken(ctx, false, resToken)
		if err != nil {
			return nil, err
		}
		if resToken.AccessToken == "" {
			return nil, fmt.Errorf("get access token error")
		}

		arrayReturn := &object.StringMap{
			key: resToken.AccessToken,
		}

		return arrayReturn, err
	}
}

func (tHandler *TokenHandler) ApplyToRequest(request *http.Request) (*http.Request, error) {
	// query Access Token power
	mapToken, err := tHandler.GetTokenQuery(request.Context())
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
