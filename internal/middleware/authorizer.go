package middleware

import (
	"api-gateway/internal/services/auth"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const userIDkey userIDContextKey = "userID"

type userIDContextKey string

func AuthorizerMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				http.Error(w, "unathorized", http.StatusUnauthorized)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := authService.Authorize(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userIDkey, claims.(jwt.MapClaims)[string(userIDkey)])
			nextHandler.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	val := ctx.Value(userIDkey)
	if val == nil {
		return 0, false
	}
	return int(val.(float64)), true
}
