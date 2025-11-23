package middleware

import (
	"api-gateway/internal/services/rate_limiter"
	"net/http"
	"strconv"
)

func LimiterMiddleware(limiterService *rate_limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			discriminator, ok := UserIDFromContext(r.Context())
			if !ok {
				http.Error(w, "unathorized", http.StatusUnauthorized)
				return
			}
			if !limiterService.Allow(strconv.Itoa(discriminator)) {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			nextHandler.ServeHTTP(w, r)
		})

	}
}
