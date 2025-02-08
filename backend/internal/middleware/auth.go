// backend/internal/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware проверяет JWT токен и, в случае успешной валидации,
// сохраняет идентификатор пользователя в контексте запроса.
func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Отсутствует токен авторизации", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Проверка корректности метода подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenInvalidId
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
				return
			}

			// Извлечение user_id из claims токена
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userIDFloat, ok := claims["user_id"].(float64)
				if !ok {
					http.Error(w, "Неверные данные токена", http.StatusUnauthorized)
					return
				}
				userID := int(userIDFloat)
				// Сохранение user_id в контексте запроса
				ctx := context.WithValue(r.Context(), "user_id", userID)
				r = r.WithContext(ctx)
			} else {
				http.Error(w, "Неверные данные токена", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
