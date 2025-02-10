// backend/internal/models/user.go
package models

// User представляет модель пользователя в системе.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // Хранит хэшированное значение пароля
}
