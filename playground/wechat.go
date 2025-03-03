package playground

import (
	"context"
	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/utils/fmt"
)

func PlayWechatOfficialAccount(localConfig *config.LocalConfig, mediaX *client.MediaX) {
	wechatOAClient, err := mediaX.CreateWechatOfficialAccount(localConfig.WeChatOfficialAccountConfig)
	if err != nil {
		panic(err)
	}

	// 调用 WeChatClient 的方法
	ctx := context.Background()
	publisher := wechatOAClient.GetPublishClient()
	res, err := publisher.PublishGet(ctx, 1)
	if err != nil {
		panic(err)
	}
	fmt.Dump(res)

	ips, err := wechatOAClient.GetCallbackIp(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Dump(ips)
}
