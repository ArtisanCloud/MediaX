package main

import (
	"github.com/ArtisanCloud/MediaX/pkg/plugin"
	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/logger"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
)

func main() {
	// init logger
	l, err := logger.NewLogger(nil, &object.HashMap{
		"level":      "info",
		"env":        "develop",
		"outputPath": "./logs/info.log",
		"errorPath":  "./logs/error.log",
		"stdout":     false,
	})
	if err != nil {
		panic(err)
	}

	l.Info("Hello, MediaX!")

	// plugin managers
	pluginManager := plugin.NewPluginManager()
	err = pluginManager.LoadPlugins()
	if err != nil {
		l.Error("Failed to load plugins:", err)
		return
	}

	fmt.Dump(pluginManager.Plugins)

}
