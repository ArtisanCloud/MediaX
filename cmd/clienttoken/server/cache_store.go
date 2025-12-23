package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type tokenCacheRecord struct {
	AppID       string    `json:"app_id"`
	AccessToken string    `json:"access_token"`
	StoredAt    time.Time `json:"stored_at"`
	ExpireAt    time.Time `json:"expire_at"`
	Source      string    `json:"source"`
	Storage     string    `json:"storage,omitempty"`
	ErrorCode   int       `json:"error_code,omitempty"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
}

func (r *tokenCacheRecord) ttl(now time.Time) time.Duration {
	if r == nil {
		return 0
	}
	return r.ExpireAt.Sub(now)
}

func (r *tokenCacheRecord) TTLSeconds(now time.Time) float64 {
	return r.ttl(now).Seconds()
}

type tokenCacheStore struct {
	redis  redis.Cmdable
	memory map[string]*tokenCacheRecord
	mu     sync.RWMutex
}

func newTokenCacheStore(redis redis.Cmdable) *tokenCacheStore {
	return &tokenCacheStore{
		redis:  redis,
		memory: make(map[string]*tokenCacheRecord),
	}
}

func (s *tokenCacheStore) save(ctx context.Context, key string, record *tokenCacheRecord, ttl time.Duration) error {
	if record == nil {
		return fmt.Errorf("token cache save: record is nil")
	}
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("token cache save: key is empty")
	}
	copyRecord := cloneTokenCacheRecord(record, "memory")
	s.mu.Lock()
	s.memory[key] = copyRecord
	s.mu.Unlock()
	if s.redis == nil {
		return nil
	}
	redisRecord := cloneTokenCacheRecord(record, "redis")
	data, err := json.Marshal(redisRecord)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, data, ttl).Err()
}

func (s *tokenCacheStore) fetch(ctx context.Context, key string) (*tokenCacheRecord, string, time.Duration, error) {
	if strings.TrimSpace(key) == "" {
		return nil, "", 0, nil
	}
	now := time.Now()
	s.mu.RLock()
	rec, ok := s.memory[key]
	s.mu.RUnlock()
	if ok && rec != nil && rec.ExpireAt.After(now) {
		record := cloneTokenCacheRecord(rec, "memory")
		record.Source = firstOrDash(record.Source)
		return record, "memory", maxDuration(record.ttl(now), 0), nil
	}
	if s.redis == nil {
		return nil, "", 0, nil
	}
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, "", 0, nil
		}
		return nil, "", 0, err
	}
	var record tokenCacheRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, "", 0, err
	}
	ttlValue, err := s.redis.TTL(ctx, key).Result()
	if err != nil || ttlValue <= 0 {
		ttlValue = record.ttl(now)
	}
	record.Storage = "redis"
	record.Source = firstOrDash(record.Source)
	s.mu.Lock()
	s.memory[key] = cloneTokenCacheRecord(&record, "memory")
	s.mu.Unlock()
	return &record, "redis", maxDuration(ttlValue, 0), nil
}

func (s *tokenCacheStore) delete(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return nil
	}
	s.mu.Lock()
	delete(s.memory, key)
	s.mu.Unlock()
	if s.redis == nil {
		return nil
	}
	return s.redis.Del(ctx, key).Err()
}

func cloneTokenCacheRecord(record *tokenCacheRecord, storage string) *tokenCacheRecord {
	if record == nil {
		return nil
	}
	cloned := *record
	cloned.Storage = storage
	return &cloned
}

func maxDuration(d time.Duration, min time.Duration) time.Duration {
	if d < min {
		return min
	}
	return d
}
