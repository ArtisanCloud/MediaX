package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/ArtisanCloud/MediaX/pkg/client"
	config2 "github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaX/playground"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	loggerconfig "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
	"github.com/ArtisanCloud/MediaXCore/utils/fmt"
)

func main() {
	configPath := resolveConfigPath()
	localConfig := mustLoadConfig(configPath)

	cacheStore, cleanup := buildCacheStore()
	defer cleanup()

	mediaX := client.NewMediaX(&config2.MediaXConfig{
		Logger: buildLoggerConfig(),
	}, cacheStore)

	if enableGoogleYouTubePlayground() {
		fmt.Dump("playground: start google youtube sample ▶️")
		playground.PlayGoogleYouTube(localConfig, mediaX)
		return
	}

	fmt.Dump("playground: no sample enabled. export PLAYGROUND_GOOGLE_YOUTUBE=1 to run Google demo.")
}

func resolveConfigPath() string {
	if env := strings.TrimSpace(os.Getenv("MEDIA_X_CONFIG")); env != "" {
		return env
	}
	return "config.yaml"
}

func mustLoadConfig(path string) *config2.LocalConfig {
	localConfig := &config2.LocalConfig{}
	if err := utils.LoadYAML(path, localConfig); err != nil {
		log.Fatalf("load config %s: %v", path, err)
	}
	return localConfig
}

func buildCacheStore() (cache.ICache, func()) {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("PLAYGROUND_CACHE_MODE")), "memory") {
		return cache.NewMemoryCache(), func() {}
	}

	addr := strings.TrimSpace(os.Getenv("PLAYGROUND_REDIS_ADDR"))
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	db := 0
	if val := strings.TrimSpace(os.Getenv("PLAYGROUND_REDIS_DB")); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			db = parsed
		}
	}
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	return cache.NewRedisCache(client), func() {
		_ = client.Close()
	}
}

func buildLoggerConfig() *loggerconfig.LogConfig {
	level := strings.TrimSpace(os.Getenv("MEDIA_X_LOG_LEVEL"))
	if level == "" {
		level = "info"
	}
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Printf("playground: create logs dir failed: %v", err)
	}
	return &loggerconfig.LogConfig{
		Level:   level,
		Console: true,
		File: loggerconfig.FileConfig{
			Enable:        true,
			InfoFilePath:  "logs/info.log",
			ErrorFilePath: "logs/error.log",
		},
	}
}

func enableGoogleYouTubePlayground() bool {
	val := strings.TrimSpace(os.Getenv("PLAYGROUND_GOOGLE_YOUTUBE"))
	return val == "1" || strings.EqualFold(val, "true")
}
