package main

import (
	"context"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/plugin/core/contract"

	"github.com/ArtisanCloud/MediaX/pkg/plugin"
	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
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

	fmt.Dump(pluginManager.Plugins)

	// 测试PluginMediaX插件
	mediaXPlugin, err := pluginManager.GetPlugin("PluginMediaX", "appId123", configPlugin)
	if err != nil {
		pluginManager.Logger.Error("Plugin MediaX Wechat not found")
		return
	}

	pluginManager.Logger.InfoF("plugin loaded name :%s \n", mediaXPlugin.Name(&ctx))

	// 插件发布内容
	resObj, err := mediaXPlugin.Publish(&ctx, &contract.PublishRequest{
		Content: "Hello, MediaX Plugin!",
	})

	if err != nil {
		panic(err)
	}
	if obj := resObj.(*contract.PublishResponse); obj.Code == 0 {
		pluginManager.Logger.Info("Publishing MediaX Plugin")
	} else {
		pluginManager.Logger.Info("Failed to publish MediaX Plugin")
	}

	fmt.Dump(pluginManager.Plugins)

	// 测试PluginMediaXWechat插件
	wechatPlugin, err := pluginManager.GetPlugin("PluginMediaXWechat", "appId456", configPlugin)
	if err != nil {
		pluginManager.Logger.Error("Plugin MediaX Wechat not found")
		return
	}

	pluginManager.Logger.InfoF("plugin loaded name :%s \n", wechatPlugin.Name(&ctx))

	// 插件发布内容
	resObj, err = wechatPlugin.Publish(&ctx, &contract.PublishRequest{
		Content: "Hello, MediaX Plugin!",
	})

	if err != nil {
		panic(err)
	}
	if obj := resObj.(*contract.PublishResponse); obj.Code == 0 {
		pluginManager.Logger.Info("Publishing MediaX Plugin")
	} else {
		pluginManager.Logger.Info("Failed to publish MediaX Plugin")
	}

	fmt.Dump(pluginManager.Plugins)
}
