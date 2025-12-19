package handler

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	jsonFieldPattern    = regexp.MustCompile(`(?i)("(?:code|state|access_token|refresh_token|client_secret)"\s*:\s*)"([^"]*)"`)
	jsonQueryPattern    = regexp.MustCompile(`(?i)((?:code|state|access_token|refresh_token|client_secret)=)([^&"]+)`)
	headerMaskKeyLookup = map[string]struct{}{
		"authorization":       {},
		"proxy-authorization": {},
		"cookie":              {},
		"set-cookie":          {},
		"x-api-token":         {},
		"x-access-token":      {},
		"x-auth-token":        {},
	}
)

// MaskSensitiveQuery hides sensitive query parameters such as code/state/token.
func MaskSensitiveQuery(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return raw
	}
	values, err := url.ParseQuery(raw)
	if err != nil {
		return raw
	}
	for _, key := range []string{"code", "state", "client_secret", "access_token", "refresh_token"} {
		if values.Has(key) {
			values.Set(key, "***")
		}
	}
	return values.Encode()
}

// MaskCallbackBody truncates and masks callback payloads depending on content type.
func MaskCallbackBody(raw, contentType string) string {
	if raw == "" {
		return raw
	}
	const maxLen = 4096
	if len(raw) > maxLen {
		raw = raw[:maxLen] + "...(truncated)"
	}
	lower := strings.ToLower(contentType)
	switch {
	case strings.Contains(lower, "application/x-www-form-urlencoded"):
		if masked := MaskSensitiveQuery(raw); masked != "" {
			return masked
		}
	case strings.Contains(lower, "json"):
		return MaskJSONSensitive(raw)
	}
	return MaskJSONSensitive(raw)
}

// MaskJSONSensitive masks key/value pairs in JSON documents.
func MaskJSONSensitive(body string) string {
	if strings.TrimSpace(body) == "" {
		return body
	}
	masked := jsonFieldPattern.ReplaceAllStringFunc(body, func(match string) string {
		return jsonFieldPattern.ReplaceAllString(match, `$1"***"`)
	})
	masked = jsonQueryPattern.ReplaceAllString(masked, `$1***`)
	return masked
}

// MaskHeaderValue hides secrets inside headers such as Authorization/Cookie.
func MaskHeaderValue(name, value string) string {
	if _, ok := headerMaskKeyLookup[strings.ToLower(strings.TrimSpace(name))]; !ok {
		return value
	}
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "authorization", "proxy-authorization":
		return MaskAuthorizationHeader(value)
	default:
		return "***"
	}
}

// MaskAuthorizationHeader strips token value and keeps scheme for readability.
func MaskAuthorizationHeader(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return trimmed
	}
	parts := strings.SplitN(trimmed, " ", 2)
	scheme := parts[0]
	var credential string
	if len(parts) > 1 {
		credential = strings.TrimSpace(parts[1])
	}
	if credential == "" {
		return scheme + " ***"
	}
	return scheme + " " + maskToken(credential)
}

func maskToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return "-"
	}
	if len(token) <= 4 {
		return "***"
	}
	return token[:2] + "***" + token[len(token)-2:]
}
