// backend/internal/handlers/auth.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register обрабатывает регистрацию нового пользователя.
func Register(db *sql.DB, jwtSecret, jwtExpiry string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		// Валидация пароля: минимум 8 символов (расширить проверку можно дополнительно)
		if len(creds.Password) < 8 {
			http.Error(w, "Пароль должен содержать минимум 8 символов", http.StatusBadRequest)
			return
		}

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(creds.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			http.Error(w, "Ошибка при обработке пароля", http.StatusInternalServerError)
			return
		}

		// Вставка пользователя в базу данных
		res, err := db.Exec(
			"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
			creds.Username,
			creds.Email,
			string(hashedPassword),
		)
		if err != nil {
			http.Error(w, "Ошибка регистрации пользователя", http.StatusInternalServerError)
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Ошибка получения id нового пользователя", http.StatusInternalServerError)
			return
		}

		// Генерация JWT токена
		tokenString, err := generateJWT(int(id), jwtSecret, jwtExpiry)
		if err != nil {
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})
	}
}

// Login обрабатывает аутентификацию пользователя.
func Login(db *sql.DB, jwtSecret, jwtExpiry string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		// Получение данных пользователя из базы
		var id int
		var email string
		var storedHashedPassword string
		var username string
		err := db.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", creds.Email).
			Scan(&id, &username, &email, &storedHashedPassword)
		if err != nil {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		// Сравнение пароля
		if err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(creds.Password)); err != nil {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		// Генерация JWT токена
		tokenString, err := generateJWT(id, jwtSecret, jwtExpiry)
		if err != nil {
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token":    tokenString,
			"username": username,
		})
	}
}

// refreshToken обрабатывает обновление JWT токена.
func RefreshToken(jwtSecret, jwtExpiry string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение токена из запроса
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Токен не найден", http.StatusUnauthorized)
			return
		}

		// Проверка токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Проверка срока действия токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Токен не валиден", http.StatusUnauthorized)
			return
		}

		// Генерация нового JWT токена
		userID := int(claims["user_id"].(float64))
		newTokenString, err := generateJWT(userID, jwtSecret, jwtExpiry)
		if err != nil {
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": newTokenString,
		})
	}
}

// generateJWT генерирует JWT токен с заданными сроком действия и идентификатором пользователя.
func generateJWT(userID int, jwtSecret, jwtExpiry string) (string, error) {
	duration, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(duration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
