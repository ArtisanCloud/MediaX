package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// HeaderSessionToken 定义 Zhihu API 需要的 SessionToken Header。
const HeaderSessionToken = "X-SessionToken"

type sessiontokenContextKey struct{}

// SessionTokenFromContext 从请求上下文中提取 X-SessionToken。
func SessionTokenFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	token, _ := ctx.Value(sessiontokenContextKey{}).(string)
	return token
}

// SessionTokenHeaderMiddleware 确保请求提供有效的 X-SessionToken，否则返回 401。
func SessionTokenHeaderMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken := strings.TrimSpace(r.Header.Get(HeaderSessionToken))
			if err := validateSessionToken(rawToken); err != nil {
				if log != nil {
					log.WithContext(r.Context()).WarnF(
						"zhihu_sessiontoken: invalid token path=%s reason=%s",
						r.URL.Path, err.Error(),
					)
				}
				writeSessionTokenError(w)
				return
			}
			ctx := context.WithValue(r.Context(), sessiontokenContextKey{}, rawToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateSessionToken(token string) error {
	if token == "" {
		return errors.New("session token missing")
	}
	if len(token) < 8 {
		return errors.New("session token too short")
	}
	if !strings.Contains(token, "=") {
		return errors.New("session token malformed")
	}
	return nil
}

func writeSessionTokenError(w http.ResponseWriter) {
	resp := map[string]string{
		"code":    "ZH_COOKIE_EXPIRED",
		"message": "session token missing or invalid",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(resp)
}
