package main

import (
	"context"
	"github.com/ArtisanCloud/MediaX/pkg/client"
	config2 "github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/utils/fmt"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 初始化 MediaX

	//cache := cache.NewMemoryCache()
	cacheClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	cache := cache.NewRedisCache(cacheClient)
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
			AppID:     "wx816158b5a5fb92cb",
			AppSecret: "7318e165c87ce29e623c6170d9adc256",
		},
	}
	wechatClient, err := mediaX.CreateWechat(wechatConfig)
	if err != nil {
		panic(err)
	}

	// 调用 WeChatClient 的方法
	ctx := context.Background()
	wechatClient.Publish()
	res, err := wechatClient.GetCallbackIP(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Dump(res)
}
