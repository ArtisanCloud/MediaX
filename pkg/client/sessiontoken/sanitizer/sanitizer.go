package sanitizer

import (
	"strings"

	"github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken/callback"
)

// MaskString 仅保留首尾字符，中间使用 * 替换。
func MaskString(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= 6 {
		return strings.Repeat("*", len(runes))
	}
	masked := make([]rune, len(runes))
	copy(masked, runes)
	for i := 3; i < len(runes)-3; i++ {
		masked[i] = '*'
	}
	return string(masked)
}

// MaskHeaders 返回新的 headers map，敏感值被遮挡。
func MaskHeaders(headers map[string]string, keys ...string) map[string]string {
	cloned := make(map[string]string, len(headers))
	for k, v := range headers {
		cloned[k] = v
	}
	if len(cloned) == 0 {
		return cloned
	}
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[strings.ToLower(k)] = struct{}{}
	}
	maskAll := len(keySet) == 0
	for k, v := range cloned {
		if maskAll {
			cloned[k] = MaskString(v)
			continue
		}
		if _, ok := keySet[strings.ToLower(k)]; ok {
			cloned[k] = MaskString(v)
		}
	}
	return cloned
}

// MaskCookies 会返回新的 cookies 切片，并脱敏 value。
func MaskCookies(cookies []callback.Cookie) []callback.Cookie {
	result := make([]callback.Cookie, len(cookies))
	for i, c := range cookies {
		result[i] = c
		result[i].Value = MaskString(c.Value)
	}
	return result
}

// MaskCredentialPayload 深拷贝 payload 并脱敏 token/cookie/header。
func MaskCredentialPayload(src *callback.CredentialPayload) *callback.CredentialPayload {
	if src == nil {
		return nil
	}
	copyPayload := *src
	copyPayload.SessionToken = MaskString(src.SessionToken)
	copyPayload.Cookies = MaskCookies(src.Cookies)
	if len(src.Headers) > 0 {
		copyPayload.Headers = MaskHeaders(src.Headers)
	}
	return &copyPayload
}
