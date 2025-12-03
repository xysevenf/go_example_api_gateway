package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthorize(t *testing.T) {
	secret := "test_secret"
	authService := NewAuthService(secret)

	t.Run("valid token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "123",
		})
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatal("unable to create token")
		}

		claims, err := authService.Authorize(tokenString)
		if err != nil {
			t.Errorf("token validation failed: %v", err)
		}
		if claims == nil {
			t.Error("nil claims")
		}
	})

	t.Run("bad token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "123",
		})
		tokenString, err := token.SignedString([]byte("bad_secret"))
		if err != nil {
			t.Fatal("unable to create token")
		}

		claims, err := authService.Authorize(tokenString)
		if err == nil {
			t.Error("expected error for bad token")
		}
		if claims != nil {
			t.Error("expected nil claims for bad token")
		}
	})
}
