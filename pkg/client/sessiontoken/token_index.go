package sessiontoken

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const sessionTokenIndexPrefix = "sessionToken:index:token:"

// sessionTokenIndexKey 返回 session_token 到 Flow 映射的缓存 key。
func sessionTokenIndexKey(token string) string {
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(trimmed))
	return sessionTokenIndexPrefix + hex.EncodeToString(sum[:])
}
