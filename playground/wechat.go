package playground

import (
	"context"
	"github.com/ArtisanCloud/MediaX/pkg/client"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/utils/fmt"
)

func PlayWechatOfficialAccount(localConfig *config.LocalConfig, mediaX *client.MediaX) {
	oaCfg, err := localConfig.ClientTokenProviders.ResolveWechatOfficialAccount(
		"",
		"",
		"",
	)
	if err != nil {
		panic(err)
	}
	wechatCfg := oaCfg.OfficialAccountConfig()
	wechatOAClient, err := mediaX.CreateWechatOfficialAccount(wechatCfg)
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

	ips, err := wechatOAClient.GetBaseClient().GetCallbackIp(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Dump(ips)
}
