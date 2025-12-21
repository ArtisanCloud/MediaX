package main

import (
	"fmt"
	"sort"
	"strings"
)

func (s *accessTokenServer) logServerEvent(event string, extra map[string]any) {
	if s == nil || s.logger == nil {
		return
	}
	fields := map[string]any{
		"listen_addr":     s.listenAddr,
		"storage_backend": s.storageBackend,
	}
	for k, v := range extra {
		fields[k] = v
	}
	pairs := make([]string, 0, len(fields))
	for k, v := range fields {
		pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
	}
	sort.Strings(pairs)
	s.logger.InfoF("accesstoken-server: event=%s %s", event, strings.Join(pairs, " "))
}

func (s *accessTokenServer) logFlowAction(action string, ctx *providerContext, flowID, tokenSource, detail string, extra map[string]any) {
	if s == nil {
		return
	}
	fields := map[string]any{
		"flow_id":             safeValue(flowID),
		"token_source":        safeValue(tokenSource),
		"token_source_detail": safeValue(detail),
		"storage_backend":     s.storageBackend,
	}
	if ctx != nil {
		fields["provider_code"] = safeValue(ctx.ProviderCode)
		fields["provider_app"] = safeValue(ctx.AppCode)
		fields["provider_auth_mode"] = safeValue(ctx.ModeKey)
		fields["config_path"] = safeValue(ctx.ConfigPath)
	}
	for k, v := range extra {
		fields[k] = v
	}
	s.logServerEvent(action, fields)
}

func safeValue(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "-"
	}
	return v
}
