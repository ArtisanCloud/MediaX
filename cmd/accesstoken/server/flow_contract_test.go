package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type flowListAPIResponse struct {
	Flows []flowListItem `json:"flows"`
}

type flowListItem struct {
	FlowID      string `json:"flow_id"`
	Status      string `json:"status"`
	MaskedToken string `json:"masked_token"`
}

type flowReplayAPIResponse struct {
	FlowID         string         `json:"flow_id"`
	Status         string         `json:"status"`
	StorageBackend string         `json:"storage_backend"`
	Payload        map[string]any `json:"payload"`
}

func TestHandleListFlowsReturnsMaskedRecords(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	rec := &oauthTokenRecord{
		ProviderCode: "bilbili",
		ProviderApp:  "content_center",
		AuthMode:     "default",
		ConfigPath:   server.defaultConfigPath,
		FlowID:       "oauth-flow-list",
		AccessToken:  "list-flow-token",
		ExpiresIn:    600,
		StoredAt:     time.Now().UTC(),
		Callback: &callbackRecord{
			Method: "GET",
			State:  "state-1",
			Body:   "***",
		},
	}
	if err := server.saveOAuthTokenRecord(rec); err != nil {
		t.Fatalf("saveOAuthTokenRecord: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/accesstoken/flows?provider_code=bilbili&provider_app=content_center&provider_auth_mode=default", nil)
	recorder := httptest.NewRecorder()
	server.handleListFlows(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", res.StatusCode)
	}
	var resp flowListAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Flows) != 1 {
		t.Fatalf("flows len=%d, want 1", len(resp.Flows))
	}
	item := resp.Flows[0]
	if item.FlowID != "oauth-flow-list" {
		t.Fatalf("flow_id=%s, want oauth-flow-list", item.FlowID)
	}
	if item.Status != "authorized" {
		t.Fatalf("status=%s, want authorized", item.Status)
	}
	if item.MaskedToken == "" || item.MaskedToken == "list-flow-token" {
		t.Fatalf("masked token not applied: %s", item.MaskedToken)
	}
}

func TestHandleFlowReplayReturnsPayload(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	rec := &oauthTokenRecord{
		ProviderCode: "bilbili",
		ProviderApp:  "content_center",
		AuthMode:     "default",
		ConfigPath:   server.defaultConfigPath,
		FlowID:       "oauth-flow-replay",
		AccessToken:  "replay-flow-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    900,
		StoredAt:     time.Now().UTC(),
		Callback: &callbackRecord{
			Method: "GET",
			State:  "state-2",
			Body:   "***",
		},
	}
	if err := server.saveOAuthTokenRecord(rec); err != nil {
		t.Fatalf("saveOAuthTokenRecord: %v", err)
	}

	body := []byte(`{"flow_id":"oauth-flow-replay"}`)
	req := httptest.NewRequest(http.MethodPost, "/accesstoken/flow/replay", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	server.handleFlowReplay(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", res.StatusCode)
	}
	var resp flowReplayAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.FlowID != "oauth-flow-replay" {
		t.Fatalf("flow_id=%s, want oauth-flow-replay", resp.FlowID)
	}
	if resp.Status != "authorized" {
		t.Fatalf("status=%s, want authorized", resp.Status)
	}
	if resp.Payload == nil {
		t.Fatalf("payload missing")
	}
	if got := resp.Payload["access_token"]; got != "replay-flow-token" {
		t.Fatalf("payload access_token=%v", got)
	}
	if got := resp.Payload["masked_token"]; got == "replay-flow-token" {
		t.Fatalf("masked_token not applied: %v", got)
	}
	if _, ok := resp.Payload["callback"]; !ok {
		t.Fatalf("callback missing in payload")
	}
}

func TestHandleListFlowsFiltersByProvider(t *testing.T) {
	server := newTestAccessTokenServer(t, testAccessTokenConfig)
	now := time.Now().UTC()
	bili := &oauthTokenRecord{
		ProviderCode: "bilbili",
		ProviderApp:  "content_center",
		AuthMode:     "default",
		ConfigPath:   server.defaultConfigPath,
		FlowID:       "oauth-bili",
		AccessToken:  "bili-token",
		StoredAt:     now,
	}
	redbook := &oauthTokenRecord{
		ProviderCode: "redbook_juguang",
		ProviderApp:  "juguang",
		AuthMode:     "default",
		ConfigPath:   server.defaultConfigPath,
		FlowID:       "oauth-redbook",
		AccessToken:  "redbook-token",
		StoredAt:     now.Add(time.Minute),
	}
	if err := server.saveOAuthTokenRecord(bili); err != nil {
		t.Fatalf("save bili flow: %v", err)
	}
	if err := server.saveOAuthTokenRecord(redbook); err != nil {
		t.Fatalf("save redbook flow: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/accesstoken/flows?provider_code=redbook_juguang", nil)
	recorder := httptest.NewRecorder()
	server.handleListFlows(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", res.StatusCode)
	}
	var resp flowListAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Flows) != 1 {
		t.Fatalf("flows len=%d want 1", len(resp.Flows))
	}
	if resp.Flows[0].FlowID != "oauth-redbook" {
		t.Fatalf("flow_id=%s want oauth-redbook", resp.Flows[0].FlowID)
	}
}
