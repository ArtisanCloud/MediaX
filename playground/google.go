package playground

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	videoSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/fmt"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

const (
	defaultVideoIDs       = "dQw4w9WgXcQ"
	defaultRegion         = "US"
	defaultAccessTokenTTL = time.Hour
)

// PlayGoogleYouTube 演示如何复用 config.yaml + GetOAuthToken 回调来调试视频接口。
// 根据环境变量控制 AccessToken 来源，输出的日志会脱敏 token，方便复制到工单中排查。
func PlayGoogleYouTube(localConfig *config.LocalConfig, mediaX *client.MediaX) {
	if localConfig == nil || localConfig.GoogleYouTubeConfig == nil {
		panic("playground: missing google_youtube_config in config.yaml")
	}

	token, source := resolvePlaygroundAccessToken(localConfig.GoogleYouTubeConfig)
	if strings.TrimSpace(token) == "" {
		panic("playground: provide GOOGLE_YOUTUBE_ACCESS_TOKEN or oauth.access_token before running example")
	}

	localConfig.GoogleYouTubeConfig.GetOAuthToken = func(key string, refresh bool) object.HashMap {
		if mediaX.Logger != nil {
			mediaX.Logger.InfoF(
				"playground: inject access token source=%s refresh=%t oauth_key=%s token=%s",
				source,
				refresh,
				valueOrDash(key),
				maskToken(token),
			)
		}
		return object.HashMap{
			"access_token": token,
			"expires_in":   defaultAccessTokenTTL.Seconds(),
		}
	}

	ytClient, err := mediaX.CreateGoogleYouTubeACClient(localConfig.GoogleYouTubeConfig)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := buildVideoListRequest()
	if mediaX.Logger != nil {
		mediaX.Logger.InfoF(
			"playground: videos.list part=%s ids=%s chart=%s region=%s oauth_key=%s",
			req.Part,
			valueOrDash(req.ID),
			valueOrDash(req.Chart),
			valueOrDash(req.RegionCode),
			valueOrDash(localConfig.GoogleYouTubeConfig.OauthKey),
		)
	}

	videoClient := ytClient.GetVideoClient()
	res, err := videoClient.List(ctx, req)
	if err != nil {
		if mediaX.Logger != nil {
			mediaX.Logger.ErrorF("playground: videos.list failed: %v", err)
		}
		panic(err)
	}

	fmt.Dump("playground: youtube#videoListResponse preview", res.PageInfo, len(res.Items))
	if len(res.Items) > 0 {
		fmt.Dump(res.Items[0].Snippet.Title, res.Items[0].Snippet.ChannelTitle)
	}
}

func buildVideoListRequest() *videoSchema.YouTubeVideoListReq {
	part := os.Getenv("PLAYGROUND_YOUTUBE_PART")
	if strings.TrimSpace(part) == "" {
		part = "snippet,contentDetails,statistics"
	}

	req := &videoSchema.YouTubeVideoListReq{
		Part: part,
	}

	if ids := strings.TrimSpace(os.Getenv("PLAYGROUND_YOUTUBE_VIDEO_IDS")); ids != "" {
		req.ID = ids
	} else {
		req.ID = defaultVideoIDs
	}

	if chart := strings.TrimSpace(os.Getenv("PLAYGROUND_YOUTUBE_CHART")); chart != "" {
		req.Chart = chart
	}
	if region := strings.TrimSpace(os.Getenv("PLAYGROUND_YOUTUBE_REGION")); region != "" {
		req.RegionCode = region
	} else if req.Chart != "" && req.RegionCode == "" {
		req.RegionCode = defaultRegion
	}

	return req
}

func resolvePlaygroundAccessToken(cfg *config.GoogleYouTubeConfig) (string, string) {
	if token := strings.TrimSpace(os.Getenv("GOOGLE_YOUTUBE_ACCESS_TOKEN")); token != "" {
		return token, "env"
	}
	if cfg != nil && cfg.ClientConfig != nil && cfg.OAuthConfig != nil {
		if token := strings.TrimSpace(cfg.OAuthConfig.AccessToken); token != "" {
			return token, "config"
		}
	}
	return "", ""
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

func valueOrDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}
