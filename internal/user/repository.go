package user

import (
	"database/sql"
	"fmt"
)

func GetUser(db *sql.DB, id int) (*User, error) {
	// SQL запрос
	query := `SELECT id, name, female, status_id FROM users WHERE id = $1`

	// Создаем переменную для результата
	user := &User{}

	err := db.QueryRow(query, id).Scan(
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
func GetAllUsers(db *sql.DB) ([]*User, error) {
	query := `SELECT id, name, female, status_id, rule_id FROM users`

	rows, err := db.Query(query) // Query() возвращает несколько строк
	if err != nil {
		return nil, err
	}
	defer rows.Close() // не забываем закрыть

	// Создаем слайс для результатов
	var users []*User

	// Итерируемся по всем строкам
	for rows.Next() {
		user := User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Female,
			&user.StatusID,
			&user.RuleID,
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

func CreateUser(db *sql.DB, user *User) (*User, error) {
	query := `
	INSERT INTO users (name, female, status_id, rule_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, female, status_id, rule_id`

	err := db.QueryRow(query, user.Name, user.Female, user.StatusID, user.RuleID).Scan(
		&user.ID,
		&user.Name,
		&user.Female,
		&user.StatusID,
		&user.RuleID,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUser(db *sql.DB, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(query, id)

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func UpdateUser(db *sql.DB, id int, name string, female *string, statusID int, ruleId int) (*User, error) {
	query := `
	UPDATE users
	SET name = $1, female = $2, status_id = $3, rule_id = $4
	WHERE id = $5
	RETURNING id, name, female, status_id, rule_id
	`

	user := &User{}

	err := db.QueryRow(query, name, female, statusID, ruleId, id).Scan(
		&user.ID,
		&user.Name,
		&user.Female,
		&user.StatusID,
		&user.RuleID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return user, nil
}
