package database

import (
	"database/sql"
	"fmt"
	"hello-go-backend/internal/models"
)

func GetUser(id int) (*models.User, error) {
	// SQL запрос
	query := `SELECT id, name, female, status_id FROM users WHERE id = $1`

	// Создаем переменную для результата
	user := &models.User{}

	err := DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Female,
		&user.StatusID,
	)

	// Проверяем ошибки
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return user, nil
}

// GetAllUsers получает всех пользователей из БД
func GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, name, female, status_id FROM users`

	rows, err := DB.Query(query) // Query() возвращает несколько строк
	if err != nil {
		return nil, err
	}
	defer rows.Close() // не забываем закрыть

	// Создаем слайс для результатов
	var users []*models.User

	// Итерируемся по всем строкам
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Female,
			&user.StatusID,
		)
		if err != nil {
			return nil, err
		}
		// Добавляем пользователя в слайс
		users = append(users, &user)
	}

	// Проверяем была ли ошибка при итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func CreateUser(user *models.User) (*models.User, error) {
	query := `
	INSERT INTO users (name, female, status_id)
	VALUES ($1, $2, $3)
	RETURNING id, name, female, status_id`

	createUser := &models.User{}

	err := DB.QueryRow(query, user.Name, user.Female, user.StatusID).Scan(
		&user.ID,
		&user.Name,
		&user.Female,
		&user.StatusID,
	)

	if err != nil {
		return nil, err
	}

	return createUser, nil
}

func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := DB.Exec(query, id)

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
