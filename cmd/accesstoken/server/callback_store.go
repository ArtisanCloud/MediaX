package main

import (
	"sync"
	"time"
)

type callbackRecord struct {
	Timestamp    time.Time         `json:"timestamp"`
	Method       string            `json:"method"`
	Query        string            `json:"query,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body"`
	ProviderCode string            `json:"provider_code,omitempty"`
	FlowID       string            `json:"flow_id,omitempty"`
	State        string            `json:"state,omitempty"`
}

type callbackLogStore struct {
	mu      sync.Mutex
	limit   int
	records []callbackRecord
}

func newCallbackLogStore(limit int) *callbackLogStore {
	if limit <= 0 {
		limit = 20
	}
	return &callbackLogStore{limit: limit}
}

func (s *callbackLogStore) append(record callbackRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append([]callbackRecord{record}, s.records...)
	if len(s.records) > s.limit {
		s.records = s.records[:s.limit]
	}
}

func (s *callbackLogStore) list() []callbackRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	copied := make([]callbackRecord, len(s.records))
	copy(copied, s.records)
	return copied
}

func (s *callbackLogStore) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = nil
}
