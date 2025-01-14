module github.com/ArtisanCloud/MediaX

go 1.23.1

replace github.com/ArtisanCloud/PowerLibs/v3 => ../../../../PowerWechat/PowerLibs

replace github.com/ArtisanCloud/PowerSocialite/v3 => ../../../../PowerWechat/PowerSocialite

require github.com/ArtisanCloud/PowerLibs/v3 v3.3.1

require (
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	go.opentelemetry.io/otel v1.4.0 // indirect
	go.opentelemetry.io/otel/trace v1.4.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
