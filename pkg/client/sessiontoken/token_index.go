package sessiontoken

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const sessiontokenIndexPrefix = "sessiontoken:index:token:"

// sessiontokenIndexKey 返回 session_token 到 Flow 映射的缓存 key。
func sessiontokenIndexKey(token string) string {
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(trimmed))
	return sessiontokenIndexPrefix + hex.EncodeToString(sum[:])
}
