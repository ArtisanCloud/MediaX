package kernel

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	request2 "github.com/ArtisanCloud/MediaX/internal/kernel/request"
	"github.com/ArtisanCloud/MediaX/internal/kernel/response"
	response2 "github.com/ArtisanCloud/MediaX/internal/kernel/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/contract"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/helper"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
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
	GetTokenQuery  func(ctx context.Context) (arrayQuery *object.StringMap, arrayHeader *object.StringMap, err error)

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
	if tHandler.GetCredentials == nil {
		return tHandler.CachePrefix
	}
	credentials := tHandler.GetCredentials()
	if credentials == nil {
		return tHandler.CachePrefix
	}
	keys := make([]string, 0, len(*credentials))
	for key := range *credentials {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var builder strings.Builder

	// 遍历 credentials map，拼接所有字段值
	for _, key := range keys {
		builder.WriteString((*credentials)[key])
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

func (tHandler *TokenHandler) sendRequest(ctx context.Context, credential *object.StringMap) ([]byte, error) {
	key := "json"
	if tHandler.RequestMethod == http.MethodGet {
		key = "query"
	}
	options := &object.HashMap{
		key: credential,
	}

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
	defer rs.Body.Close()
	bodyBytes, err := io.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
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
		customToken := tHandler.GetCustomToken(cacheKey, refresh)
		if customToken == nil {
			return fmt.Errorf("get access token error")
		}
		return assignTokenValue(resToken, customToken)
	}

	// get token from cache
	if !refresh {
		if exist, err := tHandler.Cache.Exists(ctx, cacheKey); err == nil && exist {
			value, err := tHandler.Cache.Get(ctx, cacheKey)
			if err == nil && len(value) > 0 {
				if err := json.Unmarshal(value, resToken); err == nil {
					return nil
				}
			}
		}
	}

	// request token from provider auth token api
	payload, err := tHandler.sendRequest(ctx, tHandler.GetCredentials())
	if err != nil {
		return err
	}
	if err := assignTokenValue(resToken, payload); err != nil {
		return err
	}
	var baseToken response.AccessTokenRes
	if err := json.Unmarshal(payload, &baseToken); err != nil {
		return err
	}
	if strings.TrimSpace(baseToken.AccessToken) == "" {
		return fmt.Errorf("access token empty: %s", strings.TrimSpace(string(payload)))
	}
	_, err = tHandler.SetToken(ctx, &baseToken, baseToken.ExpiresIn)

	return err
}

func assignTokenValue(dst interface{}, src interface{}) error {
	if dst == nil || src == nil {
		return errors.New("invalid token container")
	}
	if raw, ok := src.([]byte); ok {
		return json.Unmarshal(raw, dst)
	}
	if out, ok := dst.(*response.AccessTokenRes); ok {
		switch v := src.(type) {
		case *response.AccessTokenRes:
			*out = *v
			return nil
		case response.AccessTokenRes:
			*out = v
			return nil
		case []byte:
			return json.Unmarshal(v, out)
		}
	}
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}

func (tHandler *TokenHandler) OverrideGetTokenQuery() {
	tHandler.GetTokenQuery = func(ctx context.Context) (arrayQuery *object.StringMap, arrayHeader *object.StringMap, err error) {
		// set the current token key
		var key string
		if tHandler.QueryName != "" {
			key = tHandler.QueryName
		} else {
			key = tHandler.TokenKey
		}

		// get token string power
		resToken := &response.AccessTokenRes{}
		err = tHandler.GetToken(ctx, false, resToken)
		if err != nil {
			return nil, nil, err
		}
		if resToken.AccessToken == "" {
			return nil, nil, fmt.Errorf("get access token error")
		}

		arrayQuery = &object.StringMap{
			key: resToken.AccessToken,
		}

		return arrayQuery, nil, err
	}
}

func (tHandler *TokenHandler) ApplyToRequest(request *http.Request) (*http.Request, error) {
	// query Access Token power
	queryParams, headerParams, err := tHandler.GetTokenQuery(request.Context())
	if err != nil {
		return nil, err
	}
	// 设置 Query 参数
	if queryParams != nil {
		q := request.URL.Query()
		for key, value := range *queryParams {
			q.Set(key, value)
		}
		request.URL.RawQuery = q.Encode()
	}

	// 设置 Header 参数
	if headerParams != nil {
		for key, value := range *headerParams {
			request.Header.Set(key, value)
		}
	}

	return request, nil
}
