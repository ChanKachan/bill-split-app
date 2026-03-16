package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// GenerateCentrifugeToken используется для создания токена подключения к Centrifuge.
func GenerateCentrifugeToken(userID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":     userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(336 * time.Hour)), // Токен действителен 336 часа (2 недели)
		"publish": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
