module github.com/ArtisanCloud/MediaX

go 1.24.0

// replace github.com/ArtisanCloud/MediaXCore => ../MediaXCore

require (
	github.com/ArtisanCloud/MediaXCore v1.0.3
	github.com/redis/go-redis/v9 v9.8.0
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/time v0.14.0

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.6.0
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
