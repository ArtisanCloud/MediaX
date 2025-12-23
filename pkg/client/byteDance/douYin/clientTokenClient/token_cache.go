package clientTokenClient

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	douyinresponse "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
)

var nowFunc = time.Now

// clientTokenCacheRecord 记录抖音 client_token 缓存信息，结构体与 cmd/clienttoken/server/tokenCacheRecord 对齐，方便共享。
type clientTokenCacheRecord struct {
	AppID       string    `json:"app_id"`
	AccessToken string    `json:"access_token"`
	StoredAt    time.Time `json:"stored_at"`
	ExpireAt    time.Time `json:"expire_at"`
	Source      string    `json:"source"`
	ErrorCode   int       `json:"error_code,omitempty"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
}

func (r *clientTokenCacheRecord) ttl(now time.Time) time.Duration {
	if r == nil {
		return 0
	}
	return r.ExpireAt.Sub(now)
}

type tokenFetchFunc func(ctx context.Context) (*douyinresponse.ByteDanceAccessTokenRes, error)

type clientTokenCache struct {
	cache         cache.ICache
	key           string
	defaultTTL    int
	refreshBefore time.Duration
	appID         string

	mu     sync.Mutex
	memory *clientTokenCacheRecord
}

func newClientTokenCache(cfg *config.ByteDanceDouYinConfig, backend cache.ICache) *clientTokenCache {
	if cfg == nil {
		return nil
	}
	cfg.NormalizeClientTokenCredentials()
	key := strings.TrimSpace(cfg.EffectiveRedisKey(""))
	if key == "" {
		return nil
	}
	refresh := cfg.EffectiveRefreshBefore()
	if refresh <= 0 {
		refresh = 600
	}
	ttl := cfg.EffectiveTTLSeconds()
	if ttl <= 0 {
		ttl = 7000
	}
	appID := ""
	if cred := cfg.ClientTokenCredential(); cred != nil {
		appID = cred.ClientKeyValue()
	}
	return &clientTokenCache{
		cache:         backend,
		key:           key,
		defaultTTL:    ttl,
		refreshBefore: time.Duration(refresh) * time.Second,
		appID:         appID,
	}
}

func (c *clientTokenCache) ensure(ctx context.Context, fetch tokenFetchFunc) (*clientTokenCacheRecord, error) {
	record, err := c.load(ctx)
	if err != nil {
		return nil, err
	}
	if record != nil {
		if record.ttl(nowFunc()) > c.refreshBefore {
			return record, nil
		}
	}
	return c.forceRefresh(ctx, fetch)
}

func (c *clientTokenCache) forceRefresh(ctx context.Context, fetch tokenFetchFunc) (*clientTokenCacheRecord, error) {
	if c == nil {
		return nil, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	// double-check after locking
	if record, err := c.load(ctx); err == nil && record != nil {
		if record.ttl(nowFunc()) > c.refreshBefore {
			return record, nil
		}
	}
	return c.refresh(ctx, fetch)
}

func (c *clientTokenCache) refresh(ctx context.Context, fetch tokenFetchFunc) (*clientTokenCacheRecord, error) {
	if fetch == nil {
		return nil, nil
	}
	res, err := fetch(ctx)
	if err != nil {
		return nil, err
	}
	now := nowFunc()
	ttl := c.defaultTTL
	if res != nil && res.ExpiresIn > 0 {
		ttl = int(res.ExpiresIn)
	}
	if ttl <= 0 {
		ttl = 7000
	}
	record := &clientTokenCacheRecord{
		AppID:       c.appID,
		AccessToken: strings.TrimSpace(res.AccessToken),
		StoredAt:    now,
		ExpireAt:    now.Add(time.Duration(ttl) * time.Second),
		Source:      "douyin.client",
		ErrorCode:   res.ErrCode,
		ErrorMsg:    strings.TrimSpace(res.ErrMsg),
	}
	if err := c.save(ctx, record, time.Duration(ttl)*time.Second); err != nil {
		return nil, err
	}
	return record, nil
}

func (c *clientTokenCache) load(ctx context.Context) (*clientTokenCacheRecord, error) {
	if c == nil {
		return nil, nil
	}
	if c.cache == nil {
		return c.memory, nil
	}
	data, err := c.cache.Get(ctx, c.key)
	if err != nil || len(data) == 0 {
		return nil, err
	}
	var record clientTokenCacheRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *clientTokenCache) save(ctx context.Context, record *clientTokenCacheRecord, ttl time.Duration) error {
	if c == nil {
		return nil
	}
	if c.cache == nil {
		c.memory = record
		return nil
	}
	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if ttl <= 0 {
		ttl = time.Duration(c.defaultTTL) * time.Second
	}
	return c.cache.Set(ctx, c.key, payload, ttl)
}
