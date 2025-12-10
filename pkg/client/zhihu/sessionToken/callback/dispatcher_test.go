package callback

import (
	"testing"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
)

func TestNewDispatcherValidatesConfig(t *testing.T) {
	cfg := &config.ZhihuSessionTokenCallbackConfig{
		CallbackSecret: "secret",
		MaxRetry:       2,
		RetryBackoff:   []int{2, 4, 8},
	}
	dispatcher, err := NewDispatcher(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dispatcher == nil {
		t.Fatalf("expected dispatcher")
	}
}
