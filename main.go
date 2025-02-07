package main

import (
	"github.com/ArtisanCloud/MediaX/pkg/client"
	config2 "github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
)

func main() {
	// 初始化 MediaX

	cache := cache.NewMemoryCache()
	mediaX := client.NewMediaX(&config2.MediaXConfig{
		&config.LogConfig{
			Level:   "debug",
			Console: true,
			File: config.FileConfig{
				Enable: true,
			},
		},
	}, cache)

	// 创建 WeChatClient
	wechatConfig := &config2.WeChatConfig{
		AppConfig: config2.AppConfig{
			BaseUri:   "https://api.weixin.qq.com",
			ProxyUri:  "",
			Timeout:   5,
			AppID:     "wx1234567890",
			AppSecret: "1234567890",
		},
	}
	wechatClient, err := mediaX.CreateWechat(wechatConfig)
	if err != nil {
		panic(err)
	}

	// 调用 WeChatClient 的方法
	wechatClient.Publish()
}
