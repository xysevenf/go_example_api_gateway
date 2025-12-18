package main

import (
	"api-gateway/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/services/auth"
	"api-gateway/internal/services/rate_limiter"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	var logger *slog.Logger
	var logLevel slog.Level
	dev := cfg.Env == "dev"
	if dev {
		logLevel = slog.LevelDebug
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	} else {
		logLevel = slog.LevelInfo
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})).With(
			slog.String("app", cfg.AppName),
		)
	}
	r.Use(httplog.RequestLogger(logger, &httplog.Options{Level: logLevel}))

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
