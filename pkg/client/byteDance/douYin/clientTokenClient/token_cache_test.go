package clientTokenClient

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	kernelresponse "github.com/ArtisanCloud/MediaX/internal/kernel/response"
	douyinresponse "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
)

func TestClientTokenCacheEnsureRefreshBeforeThreshold(t *testing.T) {
	memCache := newMemoryCache()
	cfg := &config.ByteDanceDouYinConfig{
		ClientToken: &config.ByteDanceDouYinClientTokenCredential{
			ClientKey: "client-key",
		},
		Cache: &config.ClientTokenCacheConfig{
			RedisKey:             "clientToken:douyin:client-key",
			TTLSeconds:           120,
			RefreshBeforeSeconds: 10,
		},
	}
	store := newClientTokenCache(cfg, memCache)
	if store == nil {
		t.Fatalf("expected cache store")
	}
	counter := 0
	fetch := func(ctx context.Context) (*douyinresponse.ByteDanceAccessTokenRes, error) {
		counter++
		return &douyinresponse.ByteDanceAccessTokenRes{
			AccessTokenRes: kernelresponse.AccessTokenRes{
				AccessToken: fmt.Sprintf("token-%d", counter),
				ExpiresIn:   60,
			},
		}, nil
	}

	rec, err := store.ensure(context.Background(), fetch)
	if err != nil {
		t.Fatalf("ensure failed: %v", err)
	}
	if rec.AccessToken != "token-1" {
		t.Fatalf("unexpected token %s", rec.AccessToken)
	}
	if counter != 1 {
		t.Fatalf("expected fetch once, got %d", counter)
	}

	// TTL 仍充足，不应刷新
	rec2, err := store.ensure(context.Background(), fetch)
	if err != nil {
		t.Fatalf("ensure second failed: %v", err)
	}
	if rec2.AccessToken != "token-1" {
		t.Fatalf("token should be reused got %s", rec2.AccessToken)
	}
	if counter != 1 {
		t.Fatalf("fetch should not be called again")
	}

	// 提高 refreshBefore 触发刷新
	store.refreshBefore = time.Hour
	rec3, err := store.ensure(context.Background(), fetch)
	if err != nil {
		t.Fatalf("ensure third failed: %v", err)
	}
	if rec3.AccessToken != "token-2" {
		t.Fatalf("expected refreshed token got %s", rec3.AccessToken)
	}
	if counter != 2 {
		t.Fatalf("fetch should be called twice got %d", counter)
	}
}

type memoryCache struct {
	mu   sync.Mutex
	data map[string][]byte
}

func newMemoryCache() cache.ICache {
	return &memoryCache{
		data: make(map[string][]byte),
	}
}

func (m *memoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[key], nil
}

func (m *memoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch v := value.(type) {
	case []byte:
		m.data[key] = v
	case string:
		m.data[key] = []byte(v)
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (m *memoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.data[key]
	return ok, nil
}

func (m *memoryCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (m *memoryCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (m *memoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}
