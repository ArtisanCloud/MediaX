package kernel

import (
	"context"
	"fmt"
	request2 "github.com/ArtisanCloud/MediaX/v1/internal/kernel/request"
	response2 "github.com/ArtisanCloud/MediaX/v1/internal/kernel/response"
	"github.com/ArtisanCloud/MediaX/v1/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/v1/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/contract"
	"github.com/ArtisanCloud/MediaXCore/pkg/http/helper"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
	"io"
	"net/http"
	"time"
)

type BaseClient struct {
	HttpHelper *helper.RequestHelper
	Logger     *logger.Logger
	Cache      cache.CacheInterface

	Config   *config.AppConfig
	QueryRaw bool

	TokenHandler *AccessTokenHandler

	GetMiddlewareOfAccessToken        contract.RequestMiddleware
	GetMiddlewareOfLog                func(l *logger.Logger) contract.RequestMiddleware
	GetMiddlewareOfRefreshAccessToken func(retry int) contract.RequestMiddleware
	CheckTokenNeedRefresh             func(req *http.Request, rs *http.Response, retry int) (*http.Response, error)
}

func NewBaseClient(
	cfg *config.AppConfig,
	logger *logger.Logger, cache cache.CacheInterface,
) (*BaseClient, error) {

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

	client := &BaseClient{
		Cache:      cache,
		Logger:     logger,
		HttpHelper: h,
		Config:     cfg,
	}

	// to be setup middleware here
	client.OverrideGetMiddlewares()
	client.RegisterMiddlewares()

	return client, nil
}

func (client *BaseClient) OverrideGetMiddlewares() {
	client.OverrideGetMiddlewareOfAccessToken()
	client.OverrideGetMiddlewareOfLog()
	client.OverrideGetMiddlewareOfRefreshAccessToken()
	client.OverrideCheckTokenNeedRefresh()
}

// Registers middlewares for access token, logging, etc.
func (client *BaseClient) RegisterMiddlewares() {
	// access token
	accessTokenMiddleware := client.GetMiddlewareOfAccessToken
	// log
	logMiddleware := client.GetMiddlewareOfLog

	// check invalid access token
	checkAccessTokenMiddleware := client.GetMiddlewareOfRefreshAccessToken

	client.HttpHelper.WithMiddleware(
		accessTokenMiddleware,
		logMiddleware(client.Logger),
		checkAccessTokenMiddleware(3),
		helper.HttpDebugMiddleware(client.Config.HttpDebug),
	)
}

// OverrideGetMiddlewareOfAccessToken Middleware for access token
func (client *BaseClient) OverrideGetMiddlewareOfAccessToken() {
	client.GetMiddlewareOfAccessToken = func(handle contract.RequestHandle) contract.RequestHandle {
		return func(request *http.Request) (response *http.Response, err error) {

			if client.TokenHandler != nil {
				request, err = client.TokenHandler.ApplyToRequest(request)
			}
			if err != nil {
				return nil, err
			}

			return handle(request)
		}
	}
}

// OverrideGetMiddlewareOfLog Middleware for logging
func (client *BaseClient) OverrideGetMiddlewareOfLog() {
	client.GetMiddlewareOfLog = func(l *logger.Logger) contract.RequestMiddleware {
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

// OverrideGetMiddlewareOfRefreshAccessToken Middleware for refreshing access token
func (client *BaseClient) OverrideGetMiddlewareOfRefreshAccessToken() {
	client.GetMiddlewareOfRefreshAccessToken = func(retry int) contract.RequestMiddleware {
		return func(handle contract.RequestHandle) contract.RequestHandle {
			return func(request *http.Request) (response *http.Response, err error) {
				response, err = handle(request)
				if err != nil {
					return nil, err
				}

				if response.StatusCode != http.StatusOK {
					return response, fmt.Errorf("http response code:%d", response.StatusCode)
				}

				// Token refresh logic here if needed
				rs, err := client.CheckTokenNeedRefresh(request, response, retry)
				if err != nil {
					return rs, err
				} else if rs != nil {
					return rs, nil
				}

				return response, nil
			}
		}
	}
}

func (client *BaseClient) OverrideCheckTokenNeedRefresh() {
	client.CheckTokenNeedRefresh = func(req *http.Request, rs *http.Response, retry int) (*http.Response, error) {
		/*
			根据不同Vendor的规则，实现重载检查token是否正常获得
			如果不成功，则可以考虑重试获取
		*/
		return rs, nil
	}
}

func (client *BaseClient) HttpGet(ctx context.Context, url string, query *object.StringMap, outHeader interface{}, outBody interface{}) (interface{}, error) {
	return client.makeRequest(ctx, url, http.MethodGet, query, nil, outHeader, outBody)
}

func (client *BaseClient) HttpPost(ctx context.Context, url string, data interface{}, outHeader interface{}, outBody interface{}) (interface{}, error) {
	return client.makeRequest(ctx, url, http.MethodPost, nil, data, outHeader, outBody)
}

func (client *BaseClient) RequestRaw(ctx context.Context, url string, method string, options *object.HashMap, outHeader interface{}, outBody interface{}) (*http.Response, error) {
	return client.makeRequest(ctx, url, method, &object.StringMap{}, options, outHeader, outBody)
}

func (client *BaseClient) HttpUpload(ctx context.Context, url string, files *object.HashMap, form *request2.UploadForm, query interface{}, outHeader interface{}, outBody interface{}) (interface{}, error) {
	// 请求配置
	df := client.HttpHelper.Df().WithContext(ctx).Uri(url).Method(http.MethodPost)

	// 创建一个 map，用来存储内存中的文件
	var mems map[string]io.Reader
	if form != nil {
		mems = make(map[string]io.Reader)
		for _, content := range form.Contents {
			value, err := request2.ConvertFileObjectToReader(content.Value)
			if err != nil {
				return nil, err
			}
			mems[content.Name] = value
		}
	}

	// 获取上下文中的 header 信息
	headerValue := ctx.Value("headerKV")
	headerKV := &object.StringMap{}
	if headerValue != nil {
		headerKV = headerValue.(*object.StringMap)
	}

	// 设置 multipart 请求
	df.Multipart(func(multipart contract.MultipartDfInterface) {
		// 遍历文件列表，并添加文件到 multipart 中
		if files != nil {
			for name, path := range *files {
				multipart.FileByPath(name, path.(string))
			}
		}

		// 添加 header 中的额外字段
		for k, v := range *headerKV {
			multipart.FieldValue(k, v)
		}

		// 将内存中的文件添加到 multipart 中
		for k, v := range mems {
			multipart.FileMem(form.FileName, k, v)
		}
	})

	// 设置 query 参数
	if query != nil {
		queries := query.(*object.StringMap)
		if queries != nil {
			for k, v := range *queries {
				df.Query(k, v)
			}
		}
	}

	// 执行请求
	response, err := df.Request()
	if err != nil {
		return response, err
	}

	// 解析响应体
	if outBody != nil {
		err = client.HttpHelper.ParseResponseBodyContent(response, outBody)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// Simplified request handler
func (client *BaseClient) makeRequest(ctx context.Context, url, method string,
	query *object.StringMap, formData interface{},
	outHeader interface{}, outBody interface{},
) (*http.Response, error) {
	df := client.HttpHelper.Df().WithContext(ctx).Uri(url).Method(method)

	if query != nil {
		for k, v := range *query {
			df.Query(k, v)
		}
	}

	if formData != nil {
		df.Json(formData)
	}

	return client.executeRequest(df, outHeader, outBody)
}

func (client *BaseClient) makeRequestByEncodedData(ctx context.Context, url, method string, query *object.StringMap, formData interface{}, outHeader interface{}, outBody interface{}) (*http.Response, error) {
	df := client.HttpHelper.Df().WithContext(ctx).Uri(url).Method(method)

	if query != nil {
		for k, v := range *query {
			df.Query(k, v)
		}
	}

	if formData != nil {
		// Assuming that we need to handle the encoding separately
		encodedData := &utils.JsonEncoder{Data: formData}
		df.Any(encodedData)
	}

	return client.executeRequest(df, outHeader, outBody)
}

// Executes the HTTP request and processes the response
func (client *BaseClient) executeRequest(df contract.RequestDataflowInterface, outHeader interface{}, outBody interface{}) (*http.Response, error) {
	response, err := df.Request()
	if err != nil {
		return nil, err
	}

	if outBody != nil {
		err = client.HttpHelper.ParseResponseBodyContent(response, outBody)
		if err != nil {
			return nil, err
		}
	}

	if outHeader != nil {
		strHeader, err := object.JsonEncode(response.Header)
		if err != nil {
			return nil, err
		}
		err = object.JsonDecode([]byte(strHeader), outHeader)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}
