package models

// User структура для представления пользователя в базе данных
type User struct {
	ID       int     `json:"id"`        // ID пользователя
	Name     string  `json:"name"`      // Имя (обязательное)
	Female   *string `json:"female"`    // Пол (опционально, *string означает nullable)
	StatusID *int    `json:"status_id"` // ID статуса (опционально)
}
