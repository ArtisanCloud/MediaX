package main

import (
	"testing"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
)

func TestResolveListenAddrDefaults(t *testing.T) {
	t.Setenv("ACCESSTOKEN_LISTEN_ADDR", "")
	if got := resolveListenAddr("", ""); got != defaultListenAddr {
		t.Fatalf("resolveListenAddr() = %s, want %s", got, defaultListenAddr)
	}
}

func TestResolveListenAddrPortOnly(t *testing.T) {
	t.Setenv("ACCESSTOKEN_LISTEN_ADDR", "")
	got := resolveListenAddr("", "9090")
	want := defaultListenHost + ":9090"
	if got != want {
		t.Fatalf("resolveListenAddr(9090) = %s, want %s", got, want)
	}
}

func TestResolveListenAddrFullOverride(t *testing.T) {
	t.Setenv("ACCESSTOKEN_LISTEN_ADDR", "")
	got := resolveListenAddr("0.0.0.0:9999", "")
	if got != "0.0.0.0:9999" {
		t.Fatalf("resolveListenAddr(listen,port) = %s, want 0.0.0.0:9999", got)
	}
}

func TestBuildCacheStoreMemoryWithoutEnv(t *testing.T) {
	t.Setenv("ACCESSTOKEN_REDIS_ADDR", "memory")
	cacheStore, redisClient, closer, backend, err := buildCacheStore(nil)
	if err != nil {
		t.Fatalf("buildCacheStore err = %v", err)
	}
	defer closer()
	if backend != storageBackendMemory {
		t.Fatalf("backend = %s, want %s", backend, storageBackendMemory)
	}
	if redisClient != nil {
		t.Fatalf("expected redisClient nil when redis disabled")
	}
	if cacheStore == nil {
		t.Fatalf("expected cacheStore not nil")
	}
}

func TestBuildCacheStoreFallbackOnError(t *testing.T) {
	t.Setenv("ACCESSTOKEN_REDIS_ADDR", "127.0.0.1:1")
	t.Setenv("ACCESSTOKEN_REDIS_DB", "0")
	cacheStore, redisClient, closer, backend, err := buildCacheStore(nil)
	if err != nil {
		t.Fatalf("buildCacheStore err = %v", err)
	}
	defer closer()
	if backend != storageBackendMemory {
		t.Fatalf("backend = %s, want %s when redis unreachable", backend, storageBackendMemory)
	}
	if redisClient != nil {
		t.Fatalf("expected redisClient nil on fallback")
	}
	if cacheStore == nil {
		t.Fatalf("expected cacheStore not nil on fallback")
	}
}

func TestBuildCacheStoreUsesConfig(t *testing.T) {
	t.Setenv("ACCESSTOKEN_REDIS_ADDR", "")
	cfg := &config.AccessTokenRedisConfig{
		Addr: "memory",
		DB:   5,
	}
	cacheStore, redisClient, closer, backend, err := buildCacheStore(cfg)
	if err != nil {
		t.Fatalf("buildCacheStore err = %v", err)
	}
	defer closer()
	if backend != storageBackendMemory {
		t.Fatalf("backend = %s, want %s when config addr=memory", backend, storageBackendMemory)
	}
	if redisClient != nil {
		t.Fatalf("expected redisClient nil with memory config")
	}
	if cacheStore == nil {
		t.Fatalf("expected cacheStore not nil with config")
	}
}

func TestResolveFlowTTLSecondsEnv(t *testing.T) {
	t.Setenv("ACCESSTOKEN_FLOW_TTL_SECONDS", "123")
	if got := resolveFlowTTLSeconds(); got != 123 {
		t.Fatalf("resolveFlowTTLSeconds() = %d, want 123", got)
	}
	t.Setenv("ACCESSTOKEN_FLOW_TTL_SECONDS", "")
	if got := resolveFlowTTLSeconds(); got != defaultFlowTTLSeconds {
		t.Fatalf("resolveFlowTTLSeconds() default = %d, want %d", got, defaultFlowTTLSeconds)
	}
}

func TestListenHostFromAddr(t *testing.T) {
	tests := map[string]string{
		"":            "",
		":8080":       "",
		"8080":        "8080",
		"0.0.0.0:80":  "0.0.0.0",
		"127.0.0.1:7": "127.0.0.1",
	}
	for input, want := range tests {
		if got := listenHostFromAddr(input); got != want {
			t.Fatalf("listenHostFromAddr(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestValidateListenAddrWarning(t *testing.T) {
	server := &accessTokenServer{listenAddr: ":9000", logger: nil}
	if err := server.validateListenAddr(); err != nil {
		t.Fatalf("validateListenAddr err = %v", err)
	}
	if !server.listenAddrPublic {
		t.Fatalf("listenAddrPublic = false, want true for :port")
	}
}
