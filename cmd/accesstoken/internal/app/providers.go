package app

import "strings"

var providerAccessTokenEnv = map[string][]string{
	"google_youtube":    {"GOOGLE_YOUTUBE_ACCESS_TOKEN", "YOUTUBE_ACCESS_TOKEN"},
	"google_blogger":    {"GOOGLE_BLOGGER_ACCESS_TOKEN", "BLOGGER_ACCESS_TOKEN"},
	"byte_dance_douyin": {"DOUYIN_ACCESS_TOKEN", "BYTE_DANCE_ACCESS_TOKEN", "ACCESSTOKEN_DOUYIN_ACCESS_TOKEN"},
	"redbook_juguang":   {"REDBOOK_ACCESS_TOKEN", "XHS_ACCESS_TOKEN"},
	"bilbili":           {"BILIBILI_ACCESS_TOKEN"},
}

// EnvKeysForProvider 返回给定 provider_code 对应的 AccessToken 环境变量优先级
func EnvKeysForProvider(providerCode string) []string {
	code := strings.TrimSpace(providerCode)
	if code == "" {
		return nil
	}
	return providerAccessTokenEnv[code]
}
