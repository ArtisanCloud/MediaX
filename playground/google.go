package playground

import (
	"context"
	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/video/schema"
	"github.com/ArtisanCloud/MediaXCore/utils/fmt"
)

func PlayGoogleYouTube(localConfig *config.LocalConfig, mediaX *client.MediaX) {
	googleYouTubeClient, err := mediaX.CreateGoogleYouTube(localConfig.GoogleYouTubeConfig)
	if err != nil {
		panic(err)
	}

	// 调用 WeChatClient 的方法
	ctx := context.Background()
	video := googleYouTubeClient.GetVideoClient()
	res, err := video.List(ctx, &schema.YouTubeVideoListReq{})
	if err != nil {
		panic(err)
	}
	fmt.Dump(res)

}
