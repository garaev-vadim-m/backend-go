package models

// User структура для представления пользователя в базе данных
type User struct {
	ID       int     `json:"id"`        // ID пользователя
	Name     string  `json:"name"`      // Имя (обязательное)
	Female   *string `json:"female"`    // Пол (опционально, *string означает nullable)
	StatusID *int    `json:"status_id"` // ID статуса (опционально)
}

// CreateUserRequest структура для парсинга JSON запроса
type CreateUserRequest struct {
	Name     string  `json:"name"`
	Female   *string `json:"female"`
	StatusID *int    `json:"status_id"`
}

type UpdateUserRequest struct {
	Name     string  `json:"name"`
	Female   *string `json:"female"`
	StatusID *int    `json:"status_id"`
}
