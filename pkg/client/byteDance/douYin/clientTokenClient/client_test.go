package clientTokenClient

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

func TestEnsureClientToken_AutoRefreshBelowThreshold(t *testing.T) {
	memCache := newMemoryCache()
	client, callCount := newTestDouyinCTClient(t, memCache)

	base := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	current := base
	origNow := nowFunc
	nowFunc = func() time.Time { return current }
	t.Cleanup(func() { nowFunc = origNow })

	ctx := context.Background()

	rec1, err := client.EnsureClientToken(ctx, false)
	if err != nil {
		t.Fatalf("ensure first token: %v", err)
	}
	if rec1.AccessToken != "mock-token-1" {
		t.Fatalf("unexpected token1: %s", rec1.AccessToken)
	}
	if *callCount != 1 {
		t.Fatalf("expected 1 fetch, got %d", *callCount)
	}

	rec2, err := client.EnsureClientToken(ctx, false)
	if err != nil {
		t.Fatalf("ensure second token: %v", err)
	}
	if rec2.AccessToken != rec1.AccessToken {
		t.Fatalf("token should be reused, got %s", rec2.AccessToken)
	}
	if *callCount != 1 {
		t.Fatalf("fetch should not run again before threshold")
	}

	current = base.Add(55 * time.Second) // TTL 剩余 5 秒 < refresh_before(10)
	rec3, err := client.EnsureClientToken(ctx, false)
	if err != nil {
		t.Fatalf("ensure third token: %v", err)
	}
	if rec3.AccessToken != "mock-token-2" {
		t.Fatalf("expected refreshed token, got %s", rec3.AccessToken)
	}
	if *callCount != 2 {
		t.Fatalf("expected second fetch, got %d", *callCount)
	}
	if !rec3.StoredAt.Equal(current) {
		t.Fatalf("stored_at not updated after refresh")
	}
}

func TestEnsureClientToken_MemoryFallback(t *testing.T) {
	client, callCount := newTestDouyinCTClient(t, nil) // cache backend 缺失 → 内存降级

	base := time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC)
	current := base
	origNow := nowFunc
	nowFunc = func() time.Time { return current }
	t.Cleanup(func() { nowFunc = origNow })

	ctx := context.Background()

	rec1, err := client.EnsureClientToken(ctx, false)
	if err != nil {
		t.Fatalf("ensure first token: %v", err)
	}
	if *callCount != 1 {
		t.Fatalf("expected fetch count=1, got %d", *callCount)
	}

	current = current.Add(15 * time.Second)
	rec2, err := client.EnsureClientToken(ctx, false)
	if err != nil {
		t.Fatalf("ensure second token: %v", err)
	}
	if rec2.AccessToken != rec1.AccessToken {
		t.Fatalf("memory fallback should reuse token, got %s", rec2.AccessToken)
	}
	if *callCount != 1 {
		t.Fatalf("memory fallback should not refetch, got %d", *callCount)
	}

	current = base.Add(75 * time.Second) // TTL 过期 → 即使强制刷新也会触发获取
	rec3, err := client.EnsureClientToken(ctx, true)
	if err != nil {
		t.Fatalf("force refresh token: %v", err)
	}
	if rec3.AccessToken == rec1.AccessToken {
		t.Fatalf("force refresh should fetch new token")
	}
	if *callCount != 2 {
		t.Fatalf("force refresh should invoke fetch twice, got %d", *callCount)
	}
}

func newTestDouyinCTClient(t *testing.T, backend cache.ICache) (*ByteDanceDouYinCTClient, *int) {
	t.Helper()
	cfg := &config.ByteDanceDouYinConfig{
		ClientToken: &config.ByteDanceDouYinClientTokenCredential{
			ClientKey:    "client-key",
			ClientSecret: "client-secret",
		},
		Cache: &config.ClientTokenCacheConfig{
			RedisKey:             "clientToken:douyin:client-key",
			TTLSeconds:           120,
			RefreshBeforeSeconds: 10,
		},
	}
	store := newClientTokenCache(cfg, backend)
	if store == nil {
		t.Fatalf("token cache should not be nil")
	}
	handler := &core.ByteDanceTokenHandler{
		TokenHandler: &kernel.TokenHandler{
			CacheTokenKey: cfg.EffectiveRedisKey(""),
		},
	}
	count := 0
	handler.TokenHandler.GetCustomToken = func(key string, refresh bool) object.HashMap {
		count++
		return object.HashMap{
			"access_token": fmt.Sprintf("mock-token-%d", count),
			"expires_in":   60,
		}
	}
	client := &ByteDanceDouYinCTClient{
		DouYinConfig:       cfg,
		ClientTokenHandler: handler,
		tokenCache:         store,
	}
	return client, &count
}
