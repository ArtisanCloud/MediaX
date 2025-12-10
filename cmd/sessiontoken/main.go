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
	config2 "github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontokenredis "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/storage/redis"
	sessionhandler "github.com/ArtisanCloud/MediaX/pkg/server/handlers/session_token"
	"github.com/ArtisanCloud/MediaX/pkg/utils"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	loggerconfig "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
)

const (
	defaultConfigPath = "config.yaml"
	defaultListenAddr = ":7070"
)

func main() {
	configPathFlag := flag.String("config", "", "Path to config.yaml for SessionToken service")
	flag.Parse()

	if err := run(*configPathFlag); err != nil {
		log.Fatalf("sessiontoken: %v", err)
	}
}

func run(configPathFlag string) error {
	configPath := resolveConfigPath(configPathFlag)
	localConfig := &config2.LocalConfig{}
	if err := utils.LoadYAML(configPath, localConfig); err != nil {
		return fmt.Errorf("load config %s: %w", configPath, err)
	}
	if localConfig.ZhihuConfig == nil || localConfig.ZhihuConfig.SessionToken == nil {
		return errors.New("zhihu.sessiontoken: missing zhihu_config.sessionToken in config.yaml")
	}

	redisClient, err := buildRedisClient()
	if err != nil {
		return err
	}
	defer redisClient.Close()
	if err := pingRedis(redisClient); err != nil {
		return err
	}

	cacheStore := cache.NewRedisCache(redisClient)
	mediaX := client.NewMediaX(&config2.MediaXConfig{Logger: buildLogConfig()}, cacheStore)

	flowStore := sessiontokenredis.NewFlowStore(redisClient)
	manager, _, err := mediaX.CreateZhihuSessionTokenClient(localConfig.ZhihuConfig.SessionToken, flowStore)
	if err != nil {
		return fmt.Errorf("create zhihu session token client: %w", err)
	}

	apiToken := resolveAPIToken(localConfig.ZhihuConfig.SessionToken.Service.APIToken)
	listenAddr := resolveListenAddr()

	mux := http.NewServeMux()
	if err := sessionhandler.RegisterSessionTokenFlowCreateRoute(mux, manager, apiToken, mediaX.Logger); err != nil {
		return fmt.Errorf("register create route: %w", err)
	}
	if err := sessionhandler.RegisterSessionTokenFlowGetRoute(mux, manager, apiToken, mediaX.Logger); err != nil {
		return fmt.Errorf("register get route: %w", err)
	}

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			mediaX.Logger.ErrorF("sessiontoken: graceful shutdown failed: %v", err)
		}
	}()

	mediaX.Logger.InfoF("sessiontoken: server listening addr=%s api_token=%s", listenAddr, maskToken(apiToken))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("sessiontoken server: %w", err)
	}
	mediaX.Logger.InfoF("sessiontoken: server stopped")
	return nil
}

func resolveConfigPath(flagPath string) string {
	if strings.TrimSpace(flagPath) != "" {
		return flagPath
	}
	if env := os.Getenv("SESSIONTOKEN_CONFIG"); strings.TrimSpace(env) != "" {
		return env
	}
	return defaultConfigPath
}

func resolveListenAddr() string {
	if addr := firstNonEmptyEnv("SESSIONTOKEN_LISTEN_ADDR", "POWERX_SESSION_TOKEN_LISTEN_ADDR"); addr != "" {
		return addr
	}
	return defaultListenAddr
}

func resolveAPIToken(fallback string) string {
	token := firstNonEmptyEnv("SESSIONTOKEN_API_TOKEN", "POWERX_SESSION_TOKEN_API_TOKEN")
	if token != "" {
		return token
	}
	if strings.TrimSpace(fallback) != "" {
		return fallback
	}
	return "dev-session-token"
}

func buildRedisClient() (*redis.Client, error) {
	addr := firstNonEmptyEnv("SESSIONTOKEN_REDIS_ADDR", "POWERX_SESSION_TOKEN_REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	db := envToInt("SESSIONTOKEN_REDIS_DB", 0)
	if db == 0 {
		db = envToInt("POWERX_SESSION_TOKEN_REDIS_DB", 0)
	}
	password := firstNonEmptyEnv("SESSIONTOKEN_REDIS_PASSWORD", "POWERX_SESSION_TOKEN_REDIS_PASSWORD")
	username := firstNonEmptyEnv("SESSIONTOKEN_REDIS_USERNAME", "POWERX_SESSION_TOKEN_REDIS_USERNAME")
	opts := &redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	}
	return redis.NewClient(opts), nil
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
	logLevel := firstNonEmptyEnv("SESSIONTOKEN_LOG_LEVEL", "POWERX_SESSION_TOKEN_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Printf("sessiontoken: create logs dir failed: %v", err)
	}
	return &loggerconfig.LogConfig{
		Level:   logLevel,
		Console: true,
		File: loggerconfig.FileConfig{
			Enable:        true,
			InfoFilePath:  "logs/sessiontoken-info.log",
			ErrorFilePath: "logs/sessiontoken-error.log",
			MaxSize:       100,
			MaxBackups:    5,
			MaxAge:        7,
			Compress:      true,
		},
		UseJsonFormat: false,
	}
}

func firstNonEmptyEnv(keys ...string) string {
	for _, key := range keys {
		if val := strings.TrimSpace(os.Getenv(key)); val != "" {
			return val
		}
	}
	return ""
}

func envToInt(key string, fallback int) int {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
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
