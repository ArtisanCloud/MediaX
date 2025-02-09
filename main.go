package main

import (
	"context"
	"github.com/ArtisanCloud/MediaX/v1/pkg/client"
	config2 "github.com/ArtisanCloud/MediaX/v1/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/v1/pkg/utils"
	"github.com/ArtisanCloud/MediaX/v1/pkg/utils/fmt"
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

	localConfig := &config2.LocalConfig{}
	err := utils.LoadYAML("config.yaml", localConfig)
	if err != nil {
		panic(err)
	}
	//fmt.Dump(localConfig)

	// 创建 WeChatClient

	wechatOAClient, err := mediaX.CreateWechatOfficialAccount(localConfig.WeChatOfficialAccountConfig)
	if err != nil {
		panic(err)
	}

	// 调用 WeChatClient 的方法
	ctx := context.Background()
	publisher := wechatOAClient.Client.GetPublishClient()
	res, err := publisher.PublishGet(ctx, 1)
	if err != nil {
		panic(err)
	}
	fmt.Dump(res)

	ips, err := wechatOAClient.GetCallbackIP(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Dump(ips)
}
