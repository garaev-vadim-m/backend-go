package user

import (
	"hello-go-backend/internal/rules"
	userstatus "hello-go-backend/internal/user_status"
)

// User структура для представления пользователя в базе данных.
// Статус и роль приходят уже развёрнутыми объектами (JOIN), поэтому
// status_id/rule_id в ответе не дублируются — они нужны только на запись.
type User struct {
	ID     int                    `json:"id"`     // ID пользователя
	Name   string                 `json:"name"`   // Имя (обязательное)
	Female *string                `json:"female"` // Пол (опционально, *string означает nullable)
	Status *userstatus.UserStatus `json:"status"` // Статус (nil, если не назначен)
	Rule   *rules.Rule            `json:"rule"`   // Роль (nil, если не назначена)
}

// CreateUserRequest структура для парсинга JSON запроса.
// StatusID/RuleID — указатели: null означает «не назначено» (NULL в БД).
type CreateUserRequest struct {
	Name     string  `json:"name"`
	Female   *string `json:"female"`
	StatusID *int    `json:"status_id"`
	RuleID   *int    `json:"rule_id"`
}

type UpdateUserRequest struct {
	Name     string  `json:"name"`
	Female   *string `json:"female"`
	StatusID *int    `json:"status_id"`
	RuleID   *int    `json:"rule_id"`
}
