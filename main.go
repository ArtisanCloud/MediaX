package main

import (
	"github.com/ArtisanCloud/MediaX/pkg/client"
	config2 "github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaX/playground"
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
	c := cache.NewRedisCache(cacheClient)
	mediaX := client.NewMediaX(&config2.MediaXConfig{
		&config.LogConfig{
			Level:   "debug",
			Console: true,
			File: config.FileConfig{
				Enable: true,
			},
		},
	}, c)

	localConfig := &config2.LocalConfig{}
	err := utils.LoadYAML("config.yaml", localConfig)
	if err != nil {
		panic(err)
	}
	//fmt.Dump(localConfig)

	// 创建 WeChatClient
	playground.PlayWechatOfficialAccount(localConfig, mediaX)
	//playground.PlayGoogleYouTube(localConfig, mediaX)
}
