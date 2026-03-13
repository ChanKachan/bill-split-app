package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// GenerateCentrifugeToken используется для создания токена подключения к Centrifuge.
func GenerateCentrifugeToken(userID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":     userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен действителен 24 часа
		"publish": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
