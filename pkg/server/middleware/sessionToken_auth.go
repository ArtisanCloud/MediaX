package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

var (
	// ErrInvalidBearer 表示授权头不存在或格式错误。
	ErrInvalidBearer = errors.New("sessiontoken: invalid authorization header")
)

// SessionTokenAuthMiddleware 校验 Authorization: Bearer 头是否匹配配置。
func SessionTokenAuthMiddleware(expectedToken string, log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r.Header.Get("Authorization"))
			if token == "" || strings.TrimSpace(expectedToken) == "" || token != strings.TrimSpace(expectedToken) {
				if log != nil {
					log.WithContext(r.Context()).WarnF("sessiontoken_auth: unauthorized request, path=%s", r.URL.Path)
				}
				http.Error(w, ErrInvalidBearer.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
