package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	loggerconfig "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
)

const (
	defaultConfigPath = "config.yaml"
	defaultListenAddr = ":7072"
)

func main() {
	configFlag := flag.String("config", "", "Path to config.yaml for ClientToken service")
	portFlag := flag.String("port", "", "Listen port, e.g. 7072 or :7072")
	flag.Parse()

	if err := run(*configFlag, *portFlag); err != nil {
		log.Fatalf("clienttoken: %v", err)
	}
}

func run(configPathFlag, portFlag string) error {
	configPath := resolveConfigPath(configPathFlag)
	localCfg := &config.LocalConfig{}
	if err := utils.LoadYAML(configPath, localCfg); err != nil {
		return fmt.Errorf("load config %s: %w", configPath, err)
	}
	if localCfg.ClientTokenProviders == nil {
		return errors.New("client_token_providers 未配置")
	}

	redisClient, err := buildRedisClient()
	if err != nil {
		return err
	}
	if err := pingRedis(redisClient); err != nil {
		return err
	}

	cacheStore := cache.NewRedisCache(redisClient)
	mediaX := client.NewMediaX(&config.MediaXConfig{
		Logger: buildLogConfig(),
	}, cacheStore)

	server, err := newClientTokenServer(&clientTokenServerOptions{
		localConfig: localCfg,
		configPath:  configPath,
		listenAddr:  resolveListenAddr(portFlag),
		redis:       redisClient,
		mediaX:      mediaX,
		cache:       cacheStore,
	})
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := &http.Server{
		Addr:         server.listenAddr,
		Handler:      server.routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.logger.ErrorF("clienttoken: graceful shutdown failed: %v", err)
		}
	}()

	server.logger.InfoF("clienttoken: server listening addr=%s api_token=%s config=%s",
		server.listenAddr, maskToken(server.apiToken), configPath)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("clienttoken server: %w", err)
	}
	server.logger.InfoF("clienttoken: server stopped")
	return nil
}

func resolveConfigPath(flagPath string) string {
	if val := strings.TrimSpace(flagPath); val != "" {
		return val
	}
	if env := strings.TrimSpace(os.Getenv("CLIENTTOKEN_CONFIG")); env != "" {
		return env
	}
	return defaultConfigPath
}

func resolveListenAddr(flagPort string) string {
	port := strings.TrimSpace(flagPort)
	if port != "" {
		if strings.HasPrefix(port, ":") {
			return port
		}
		return ":" + port
	}
	if env := strings.TrimSpace(os.Getenv("CLIENTTOKEN_LISTEN_ADDR")); env != "" {
		return env
	}
	return defaultListenAddr
}

func buildRedisClient() (*redis.Client, error) {
	addr := strings.TrimSpace(os.Getenv("CLIENTTOKEN_REDIS_ADDR"))
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	db := 0
	if val := strings.TrimSpace(os.Getenv("CLIENTTOKEN_REDIS_DB")); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			db = parsed
		}
	}
	username := strings.TrimSpace(os.Getenv("CLIENTTOKEN_REDIS_USERNAME"))
	password := strings.TrimSpace(os.Getenv("CLIENTTOKEN_REDIS_PASSWORD"))
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Username: username,
		Password: password,
	})
	return client, nil
}

func pingRedis(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}
	return nil
}

func buildLogConfig() *loggerconfig.LogConfig {
	level := strings.TrimSpace(os.Getenv("CLIENTTOKEN_LOG_LEVEL"))
	if level == "" {
		level = "info"
	}
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Printf("clienttoken: mkdir logs failed: %v", err)
	}
	return &loggerconfig.LogConfig{
		Level:   level,
		Console: true,
		File: loggerconfig.FileConfig{
			Enable:        true,
			InfoFilePath:  "logs/clienttoken-server-info.log",
			ErrorFilePath: "logs/clienttoken-server-error.log",
			MaxSize:       50,
			MaxBackups:    3,
			MaxAge:        7,
		},
	}
}

func maskToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return "-"
	}
	if len(token) <= 4 {
		return "***"
	}
	return token[:2] + "***" + token[len(token)-2:]
}
