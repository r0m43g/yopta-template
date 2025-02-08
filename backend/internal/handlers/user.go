// backend/internal/handlers/user.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserProfile struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// GetProfile возвращает профиль пользователя.
func GetProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromContext(r)
		if err != nil {
			http.Error(w, "Не удалось получить идентификатор пользователя", http.StatusUnauthorized)
			return
		}

		var profile UserProfile
		err = db.QueryRow("SELECT id, username FROM users WHERE id = ?", userID).
			Scan(&profile.ID, &profile.Username)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

type PasswordChangeRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword позволяет сменить пароль пользователя, предварительно проверяя старый.
func ChangePassword(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromContext(r)
		if err != nil {
			http.Error(w, "Не удалось получить идентификатор пользователя", http.StatusUnauthorized)
			return
		}

		var req PasswordChangeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		// Получение текущего хэшированного пароля из базы
		var storedHashedPassword string
		err = db.QueryRow("SELECT password FROM users WHERE id = ?", userID).
			Scan(&storedHashedPassword)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		// Проверка старого пароля
		if err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(req.OldPassword)); err != nil {
			http.Error(w, "Неверный старый пароль", http.StatusUnauthorized)
			return
		}

		// Валидация нового пароля
		if len(req.NewPassword) < 8 {
			http.Error(w, "Новый пароль должен содержать минимум 8 символов", http.StatusBadRequest)
			return
		}

		// Хэширование нового пароля
		newHashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(req.NewPassword),
			bcrypt.DefaultCost,
		)
		if err != nil {
			http.Error(w, "Ошибка обработки нового пароля", http.StatusInternalServerError)
			return
		}

		// Обновление пароля в базе данных
		_, err = db.Exec(
			"UPDATE users SET password = ? WHERE id = ?",
			string(newHashedPassword),
			userID,
		)
		if err != nil {
			http.Error(w, "Ошибка обновления пароля", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пароль успешно обновлен"))
	}
}

// getUserIDFromContext извлекает идентификатор пользователя из контекста запроса.
func getUserIDFromContext(r *http.Request) (int, error) {
	// Идентификатор пользователя устанавливается JWT middleware
	userIDValue := r.Context().Value("user_id")
	if userIDValue == nil {
		return 0, http.ErrNoCookie
	}
	userID, ok := userIDValue.(int)
	if !ok {
		return 0, http.ErrNoCookie
	}
	return userID, nil
}
