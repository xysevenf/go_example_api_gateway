package main

import (
	"api-gateway/config"
	"api-gateway/internal/services/auth"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestConfig() *config.Config {
	cfg := &config.Config{}

	cfg.Auth.JwtSignToken = "secret123"

	cfg.Limiter.MaxReq = 2
	cfg.Limiter.Window = 60

	return cfg
}

func newAuthenticatedRequest(method, path, secret string, userID int) *http.Request {
	authService := auth.NewAuthService(secret)
	token, _ := authService.IssueJWT(userID)

	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func TestRouter_AuthorizedRequest(t *testing.T) {
	cfg := newTestConfig()
	r := NewRouter(cfg)

	req := newAuthenticatedRequest("GET", "/", cfg.Auth.JwtSignToken, 123)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestRouter_AuthMissing(t *testing.T) {
	cfg := newTestConfig()
	r := NewRouter(cfg)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestRouter_InvalidToken(t *testing.T) {
	cfg := newTestConfig()
	r := NewRouter(cfg)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer BADTOKEN")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestRouter_RateLimitExceeded(t *testing.T) {
	cfg := newTestConfig()
	cfg.Limiter.MaxReq = 1

	r := NewRouter(cfg)

	req := newAuthenticatedRequest("GET", "/", cfg.Auth.JwtSignToken, 111)

	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req)

	if w1.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w1.Code)
	}

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)

	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", w2.Code)
	}
}
