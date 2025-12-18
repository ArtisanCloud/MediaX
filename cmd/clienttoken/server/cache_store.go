package main

import (
	"context"
	"encoding/json"
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
	ErrorCode   int       `json:"error_code,omitempty"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
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
	s.mu.Lock()
	s.memory[key] = record
	s.mu.Unlock()
	if s.redis == nil {
		return nil
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, data, ttl).Err()
}

func (s *tokenCacheStore) fetch(ctx context.Context, key string) (*tokenCacheRecord, string, error) {
	s.mu.RLock()
	rec, ok := s.memory[key]
	s.mu.RUnlock()
	if ok && rec != nil && rec.ExpireAt.After(time.Now()) {
		rec.Source = "memory"
		return rec, "memory", nil
	}
	if s.redis == nil || key == "" {
		return nil, "", nil
	}
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, "", nil
		}
		return nil, "", err
	}

	var record tokenCacheRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, "", err
	}
	s.mu.Lock()
	s.memory[key] = &record
	s.mu.Unlock()
	record.Source = "redis"
	return &record, "redis", nil
}
