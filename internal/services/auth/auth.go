package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	secret string
}

func (a *AuthService) Authorize(token string) (jwt.Claims, error) {
	parsedToken, err := jwt.Parse(token, func(_ *jwt.Token) (any, error) {
		return []byte(a.secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("bad token")
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("bad token")
	}
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{secret: secret}
}
