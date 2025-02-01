module github.com/ArtisanCloud/MediaX

go 1.23

toolchain go1.23.1

replace github.com/ArtisanCloud/PowerSocialite/v3 => ../../../../PowerWechat/PowerSocialite

require (
	github.com/ArtisanCloud/MediaXCore v0.0.0-20250201094827-38ed39ee85fa
	github.com/ArtisanCloud/MediaXPlugin v0.0.0-20250201122009-d5c37a6d9379
	github.com/ArtisanCloud/MediaXPlugin-Wechat v0.0.0-20250201122858-b54edc04037e
	github.com/ArtisanCloud/PowerLibs/v3 v3.3.1
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
