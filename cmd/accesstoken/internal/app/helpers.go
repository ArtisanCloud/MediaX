package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	loggerconfig "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
)

// ResolveConfigPath picks config path from flag/env/default.
func ResolveConfigPath(flagPath string) string {
	if val := strings.TrimSpace(flagPath); val != "" {
		return val
	}
	if env := FirstNonEmptyEnv("ACCESSTOKEN_CONFIG", "MEDIA_X_CONFIG"); env != "" {
		return env
	}
	return DefaultConfigPath
}

// ResolveAccessToken applies CLI/source priority: flag > env > config.
func ResolveAccessToken(flagToken, providerCode string, clientCfg *config.ClientConfig) (string, string) {
	if token := strings.TrimSpace(flagToken); token != "" {
		return token, "flag"
	}
	if token, key := firstEnvToken(EnvKeysForProvider(providerCode)); token != "" {
		if key == "" {
			return token, "env"
		}
		return token, "env:" + key
	}
	if clientCfg != nil && clientCfg.OAuthConfig != nil {
		if token := strings.TrimSpace(clientCfg.OAuthConfig.AccessToken); token != "" {
			return token, "config"
		}
	}
	return "", ""
}

// FirstNonEmptyEnv returns the first non-empty environment variable.
func FirstNonEmptyEnv(keys ...string) string {
	for _, k := range keys {
		if val := strings.TrimSpace(os.Getenv(k)); val != "" {
			return val
		}
	}
	return ""
}

func firstEnvToken(keys []string) (string, string) {
	for _, k := range keys {
		if val := strings.TrimSpace(os.Getenv(k)); val != "" {
			return val, k
		}
	}
	return "", ""
}

// BuildBasicLogConfig builds a simple console logger configuration.
func BuildBasicLogConfig(level string) *loggerconfig.LogConfig {
	if strings.TrimSpace(level) == "" {
		level = "info"
	}
	return &loggerconfig.LogConfig{
		Level:   level,
		Console: true,
	}
}

// BuildFileLogConfig extends BuildBasicLogConfig with file outputs.
func BuildFileLogConfig(level, infoFile, errorFile string) *loggerconfig.LogConfig {
	cfg := BuildBasicLogConfig(level)
	if err := os.MkdirAll("logs", 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "accesstoken: mkdir logs failed: %v\n", err)
	}
	cfg.File = loggerconfig.FileConfig{
		Enable:        true,
		InfoFilePath:  infoFile,
		ErrorFilePath: errorFile,
		MaxSize:       100,
		MaxBackups:    5,
		MaxAge:        7,
		Compress:      false,
	}
	return cfg
}

// LogInvocation prints a structured log before executing the request.
func LogInvocation(log *logger.Logger, opts *Options, oauthKey string) {
	if log == nil || opts == nil {
		return
	}
	providerCode := ValueOrDash(opts.ProviderCode)
	providerApp := ValueOrDash(opts.ProviderApp)
	log.InfoF(
		"accesstoken: provider=%s app=%s action=%s part=%s ids=%s channel_id=%s query=%s mine=%t search_mine=%t search_type=%s max_results=%d region=%s page_token=%s chart=%s category=%s token_source=%s token=%s oauth_key=%s",
		providerCode,
		providerApp,
		opts.Action,
		ValueOrDash(opts.Part),
		ValueOrDash(opts.IDs),
		ValueOrDash(opts.ChannelID),
		ValueOrDash(opts.Query),
		opts.Mine,
		opts.SearchMine,
		ValueOrDash(opts.SearchType),
		opts.MaxResults,
		ValueOrDash(opts.Region),
		ValueOrDash(opts.PageToken),
		ValueOrDash(opts.Chart),
		ValueOrDash(opts.VideoCategory),
		ValueOrDash(opts.TokenSource),
		MaskToken(opts.AccessToken),
		ValueOrDash(oauthKey),
	)
}

// ValueOrDash keeps log formatting aligned between CLI & server.
func ValueOrDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}

// MaskToken redacts token values in logs.
func MaskToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return "-"
	}
	if len(token) <= 4 {
		return "***"
	}
	return token[:2] + "***" + token[len(token)-2:]
}
