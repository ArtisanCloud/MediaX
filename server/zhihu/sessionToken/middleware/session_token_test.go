package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSessionTokenHeaderMiddlewareRejectsMissingHeader(t *testing.T) {
	middleware := SessionTokenHeaderMiddleware(nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when header missing, got %d", rec.Code)
	}
	if got := rec.Body.String(); len(got) == 0 || !strings.Contains(got, `"code":"ZH_COOKIE_EXPIRED"`) {
		t.Fatalf("expected error payload with ZH_COOKIE_EXPIRED, got %s", got)
	}
}

func TestSessionTokenHeaderMiddlewareRejectsMalformedHeader(t *testing.T) {
	middleware := SessionTokenHeaderMiddleware(nil)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/test", nil)
	req.Header.Set(HeaderSessionToken, "invalid-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for malformed token, got %d", rec.Code)
	}
}

func TestSessionTokenHeaderMiddlewarePassesValidToken(t *testing.T) {
	middleware := SessionTokenHeaderMiddleware(nil)
	captured := ""
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = SessionTokenFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/test", nil)
	req.Header.Set(HeaderSessionToken, "SESSIONID=abc123; JOID=xyz")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if captured == "" {
		t.Fatalf("expected session token to be stored in context")
	}
}
