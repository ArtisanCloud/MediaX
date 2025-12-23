package config

import (
	"testing"
)

func TestResolveByteDanceDouYinClientToken(t *testing.T) {
	cfg := &ClientTokenProvidersConfig{
		Providers: []*ClientTokenProvider{
			{
				Code: "byte_dance",
				Name: "字节跳动",
				Apps: []*ClientTokenProviderApp{
					{
						Code:         "douyin_service",
						ProviderCode: "byte_dance_douyin_clienttoken",
						AuthModes: []*ClientTokenAuthMode{
							{
								Key: "default",
								ByteDanceDouYinConfig: &ByteDanceDouYinConfig{
									ClientToken: &ByteDanceDouYinClientTokenCredential{
										ClientKey:    "test-key",
										ClientSecret: "test-secret",
									},
									Cache: &ClientTokenCacheConfig{
										RedisKey:             "clientToken:douyin:test-key",
										TTLSeconds:           7000,
										RefreshBeforeSeconds: 600,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	douyinCfg, err := cfg.ResolveByteDanceDouYinClientToken("byte_dance_douyin_clienttoken", "douyin_service", "default")
	if err != nil {
		t.Fatalf("expected config, got error: %v", err)
	}

	if douyinCfg == nil || douyinCfg.ClientTokenCredential().ClientKeyValue() != "test-key" {
		t.Fatalf("unexpected client token config: %+v", douyinCfg)
	}

	if douyinCfg.EffectiveRedisKey("") != "clientToken:douyin:test-key" {
		t.Fatalf("unexpected redis key: %s", douyinCfg.EffectiveRedisKey(""))
	}

	if douyinCfg.EffectiveTTLSeconds() != 7000 {
		t.Fatalf("unexpected TTL: %d", douyinCfg.EffectiveTTLSeconds())
	}

	if douyinCfg.EffectiveRefreshBefore() != 600 {
		t.Fatalf("unexpected refresh_before: %d", douyinCfg.EffectiveRefreshBefore())
	}
}
