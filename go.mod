module github.com/ArtisanCloud/MediaX/v1

go 1.18

//replace github.com/ArtisanCloud/PowerSocialite/v3 => ../../../../PowerWechat/PowerSocialite

replace github.com/ArtisanCloud/MediaXCore => ../MediaXCore

require github.com/ArtisanCloud/MediaXCore v0.0.0-20250207094210-bd092bfda3be

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
