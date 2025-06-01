package playground

import (
	"context"

	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/fmt"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

func PlayGoogleYouTube(localConfig *config.LocalConfig, mediaX *client.MediaX) {
	googleYouTubeClient, err := mediaX.CreateGoogleYouTubeACClient(localConfig.GoogleYouTubeConfig)
	if err != nil {
		panic(err)
	}

	// 设置AccessToken
	googleYouTubeClient.GoogleClient.TokenHandler.GetCustomToken = func(key string, refresh bool) object.HashMap {
		fmt.Dump("GetCustomToken", key, refresh)
		return object.HashMap{
			// 这个acess token需要开发这来维护，或者可以通过MediaX Studio的UI界面来维护
			"access_token": "72_ggzUdSgH99StJ2EhmuaIbHHUP9_3rDvdnQVQ9eoX5gwmNfuLpJgBUb5uPgdoh4aoVv9jYz3EKglRT73ppWqgRwzirNQM-bHaToDQ83ux1sFdCr5GK7jxYQfAESoCOEaAHAKWM",
			"expires_in":   float64(7200),
		}
	}

	// 调用 Youtube的VideoClient 的方法
	ctx := context.Background()
	video := googleYouTubeClient.GetVideoClient()
	res, err := video.List(ctx, &schema.YouTubeVideoListReq{})
	if err != nil {
		panic(err)
	}
	fmt.Dump(res)
}
