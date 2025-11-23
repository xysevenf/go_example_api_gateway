package main

import (
	"api-gateway/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/services/auth"
	"api-gateway/internal/services/rate_limiter"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)

	r.Use(middleware.AuthorizerMiddleware(auth.NewAuthService(cfg.Auth.JwtSignToken)))

	limiterOptions := rate_limiter.NewLimiterOptions()
	limiterOptions.MaxReq = cfg.Limiter.MaxReq
	limiterOptions.Window = cfg.Limiter.Window
	r.Use(middleware.LimiterMiddleware(rate_limiter.NewRateLimiter(limiterOptions)))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userID, _ := middleware.UserIDFromContext(r.Context())
		w.Write([]byte("welcome " + strconv.Itoa(userID)))
	})
	return r
}
