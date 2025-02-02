package main

import (
	"context"
	"fmt"
	fmt2 "github.com/ArtisanCloud/MediaX/pkg/utils/fmt"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/plugin/core/contract"

	"github.com/ArtisanCloud/MediaX/pkg/plugin"
)

func main() {
	ctx := context.Background()

	// init logger
	configPlugin := &contract.PluginConfig{
		LogConfig: config.LogConfig{
			Level:         "debug",
			Console:       true,
			UseJsonFormat: true,
			File: config.FileConfig{
				Enable: true,
				//InfoFilePath:  "./logs/info.log",
				//ErrorFilePath: "./logs/error.log",
			},
			//Loki: config.LokiConfig{},
			HttpDebug: true,
			Debug:     true,
		},
	}

	// plugin managers
	pluginManager := plugin.NewPluginManager(&configPlugin.LogConfig)
	pluginManager.Logger.Info("Starting MediaX Plugin...")
	pluginManager.Logger.Info("Hello, MediaX!")
	err := pluginManager.ScanPlugins()
	if err != nil {
		pluginManager.Logger.ErrorF("Failed to load plugins:", err.Error())
		return
	}

	fmt2.Dump(pluginManager.Plugins)

	// 测试PluginMediaXWechat插件
	err = ValidateLoadedPlugin(ctx, pluginManager, contract.WechatMediaVendor, contract.WechatOfficialAccount, "appId123", configPlugin)
	if err != nil {
		panic(err)
	}
	// 测试PluginMediaXDouYin插件
	err = ValidateLoadedPlugin(ctx, pluginManager, contract.DouYinMediaVendor, contract.DouYin, "appId456", configPlugin)
	if err != nil {
		panic(err)
	}
	// 测试PluginMediaXRedBook插件
	err = ValidateLoadedPlugin(ctx, pluginManager, contract.RedBookMediaVendor, contract.RedBook, "appId789", configPlugin)
	if err != nil {
		panic(err)
	}

}

func ValidateLoadedPlugin(
	ctx context.Context, pluginManager *plugin.PluginManager,
	vendorName contract.MediaVendor, pluginName contract.AppPlugin, appId string,
	config *contract.PluginConfig,
) error {
	loadedPlugin, err := pluginManager.GetPlugin(vendorName, pluginName, appId, config)
	if err != nil {
		pluginManager.Logger.ErrorF("load plugin %s err:%s", pluginName, err.Error())
		return err
	}
	pluginManager.Logger.InfoF("plugin loaded name :%s \n", loadedPlugin.Name(&ctx))
	// 插件发布内容
	resObj, err := loadedPlugin.Publish(&ctx, &contract.PublishRequest{
		Content: fmt.Sprintf("Hello, %s!", loadedPlugin.Name(&ctx)),
	})

	if err != nil {
		return err
	}
	if obj := resObj.(*contract.PublishResponse); obj.Code == 0 {
		pluginManager.Logger.InfoF("Published %s", loadedPlugin.Name(&ctx))
	} else {
		pluginManager.Logger.InfoF("Failed to publish %s", loadedPlugin.Name(&ctx))
	}

	fmt2.Dump(pluginManager.Plugins)
	return nil
}
