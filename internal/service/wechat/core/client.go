package core

import (
	"bytes"
	"github.com/ArtisanCloud/MediaX/v1/internal/kernel"
	response2 "github.com/ArtisanCloud/MediaX/v1/internal/kernel/response"
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/core/response"
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/officialAccount/material"
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/officialAccount/media"
	"github.com/ArtisanCloud/MediaX/v1/internal/service/wechat/officialAccount/publish"
	"github.com/ArtisanCloud/MediaX/v1/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"io"
	"net/http"
	"strings"
)

type WeChatClient struct {
	*kernel.BaseClient
	Config *config.WeChatOfficialAccountConfig

	// clients
	media    *media.Client
	material *material.Client
	publish  *publish.Client
}

func NewWeChatClient(cfg *config.WeChatOfficialAccountConfig, logger *logger.Logger, cache cache.CacheInterface) (*WeChatClient, error) {
	baseClient, err := kernel.NewBaseClient(&cfg.AppConfig, logger, cache)
	if err != nil {
		return nil, err
	}
	wechatClient := &WeChatClient{
		BaseClient: baseClient,
		Config:     cfg,
	}

	wechatClient.OverrideCheckTokenNeedRefresh()

	return wechatClient, nil
}

func (client *WeChatClient) OverrideCheckTokenNeedRefresh() {
	client.BaseClient.CheckTokenNeedRefresh = func(req *http.Request, rs *http.Response, retry int) (*http.Response, error) {
		ctx := req.Context()

		RetryDecider := func(code int) bool {

			if code == 40001 || code == 40014 || code == 42001 {
				return true
			}
			return false
		}

		// 如何微信返回的是二进制数据流，那么就无须判断返回的err code是否正常
		if client.QueryRaw {
			if !strings.Contains(rs.Header.Get("Content-Type"), "application/json") {
				return rs, nil
			}
		}

		res := &response.WeChatAccessTokenRes{}
		err := response2.ParseResponseToObject(rs, res)
		if err != nil {
			return nil, err
		}

		if res != nil && res.ErrCode > 0 {
			if retry > 0 && RetryDecider(res.ErrCode) {
				client.TokenHandler.Refresh(ctx)

				// clone 一个request
				client.Logger.WithContext(ctx).InfoF("refresh token, retry:%d", retry)
				token, err := client.TokenHandler.GetToken(ctx, false)
				if err != nil {
					return nil, err
				}
				q := req.URL.Query()
				q.Set(client.TokenHandler.TokenKey, token.AccessToken)
				req.URL.RawQuery = q.Encode()
				req2 := req.Clone(ctx)
				if req.Body != nil {
					// 缓存请求body
					reqData, err := io.ReadAll(req.Body)
					if err != nil {
						return nil, err
					}

					// 给两个request复制缓存下来的body数据
					req.Body = io.NopCloser(bytes.NewBuffer(reqData))
					req2.Body = io.NopCloser(bytes.NewReader(reqData))
				}

				res2, err := client.HttpHelper.GetClient().DoRequest(req2)
				if err != nil {
					return nil, err
				}

				return res2, err
			}
		}

		return rs, nil
	}
}

func (client *WeChatClient) GetMediaClient() *media.Client {
	if client.media == nil {
		client.media = media.NewClient(client.BaseClient)
	}
	return client.media
}

func (client *WeChatClient) GetMaterialClient() *material.Client {
	if client.material == nil {
		client.material = material.NewClient(client.BaseClient)
	}
	return client.material
}

func (client *WeChatClient) GetPublishClient() *publish.Client {
	if client.publish == nil {
		client.publish = publish.NewClient(client.BaseClient)
	}
	return client.publish
}
