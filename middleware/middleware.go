package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware создает middleware для проверки JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Пытаемся получить токен из Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}

		// 2. Если нет в header, ищем в form-data
		if tokenString == "" {
			// Из form-data или x-www-form-urlencoded
			tokenString = c.PostForm("Authorization") // Ищем поле Authorization в form-data

			// Если нашли и оно начинается с "Bearer "
			if tokenString != "" && strings.HasPrefix(tokenString, "Bearer ") {
				tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			}
		}

		// 3. Проверяем наличие токена
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "token not found in headers or form-data",
			})
			c.Abort()
			return
		}

		// Получаем секретный ключ из переменных окружения
		jwtKey := os.Getenv("jwtSecretKey")
		if jwtKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal server error",
				"message": "JWT secret key not configured",
			})
			c.Abort()
			return
		}

		// Парсим и валидируем токен
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Извлекаем userID из claims (sub)
		userIDstr, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid token claims",
			})
			c.Abort()
			return
		}

		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid token claims",
			})
			c.Abort()
			return
		}

		// Сохраняем userID в контексте
		c.Set("userID", userID)

		// Также сохраняем остальные claims если нужно
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware опциональная авторизация (не блокирует запрос если токена нет)
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		jwtKey := os.Getenv("jwtSecretKey")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err == nil && token.Valid {
			if userID, ok := claims["sub"].(string); ok {
				c.Set("userID", userID)
				c.Set("claims", claims)
			}
		}

		c.Next()
	}
}
