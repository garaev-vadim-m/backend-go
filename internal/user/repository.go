package user

import (
	"database/sql"
	"fmt"

	"hello-go-backend/internal/rules"
	userstatus "hello-go-backend/internal/user_status"
)

// Общие части SELECT для таблицы users: колонки и JOIN-ы вынесены в
// константы, чтобы list/get/create/update не разъезжались по формату.
// Алиас источника везде один и тот же — u (users, ins CTE или upd CTE).
const userColumns = `
	u.id, u.name, u.female,
	us.id, us.name, us.code,
	r.id, r.name, r.code`

const userJoins = `
	LEFT JOIN user_status us ON us.id = u.status_id
	LEFT JOIN rules r ON r.id = u.rule_id`

// scanner — общий интерфейс для *sql.Row и *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

// scanUser разбирает строку выборки (уже с JOIN-нутым статусом и ролью)
// в модель User. Nullable-колонки берём через sql.Null*, чтобы
// пользователь без статуса/роли не терялся, а получал status: null.
func scanUser(s scanner) (*User, error) {
	u := &User{}

	var statusID sql.NullInt64
	var statusName, statusCode sql.NullString
	var ruleID sql.NullInt64
	var ruleName, ruleCode sql.NullString

	err := s.Scan(
		&u.ID,
		&u.Name,
		&u.Female,
		&statusID,
		&statusName,
		&statusCode,
		&ruleID,
		&ruleName,
		&ruleCode,
	)
	if err != nil {
		return nil, err
	}

	if statusID.Valid {
		u.Status = &userstatus.UserStatus{
			ID:   int(statusID.Int64),
			Name: statusName.String,
			Code: statusCode.String,
		}
	}

	if ruleID.Valid {
		u.Rule = &rules.Rule{
			ID:   int(ruleID.Int64),
			Name: ruleName.String,
			Code: ruleCode.String,
		}
	}

	return u, nil
}

// GetUser возвращает пользователя по id вместе со статусом и ролью.
func GetUser(db *sql.DB, id int) (*User, error) {
	query := `SELECT` + userColumns + ` FROM users u` + userJoins + ` WHERE u.id = $1`

	u, err := scanUser(db.QueryRow(query, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return u, nil
}

// GetAllUsers получает всех пользователей из БД вместе со статусом и ролью.
func GetAllUsers(db *sql.DB) ([]*User, error) {
	query := `SELECT` + userColumns + ` FROM users u` + userJoins

	rows, err := db.Query(query) // Query() возвращает несколько строк
	if err != nil {
		return nil, err
	}
	defer rows.Close() // не забываем закрыть

	// Создаем слайс для результатов
	var users []*User

	// Итерируемся по всем строкам
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		// Добавляем пользователя в слайс
		users = append(users, u)
	}

	// Проверяем была ли ошибка при итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// CreateUser создаёт пользователя и возвращает его со статусом и ролью.
// statusID/ruleID — указатели: nil записывается как NULL.
func CreateUser(db *sql.DB, name string, female *string, statusID *int, ruleID *int) (*User, error) {
	// RETURNING не умеет JOIN, поэтому вставляем в CTE и джойним уже его.
	query := `
	WITH ins AS (
		INSERT INTO users (name, female, status_id, rule_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, female, status_id, rule_id
	)
	SELECT` + userColumns + ` FROM ins u` + userJoins

	u, err := scanUser(db.QueryRow(query, name, female, statusID, ruleID))
	if err != nil {
		return nil, err
	}

	return u, nil
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

// UpdateUser обновляет пользователя и возвращает его со статусом и ролью.
func UpdateUser(db *sql.DB, id int, name string, female *string, statusID *int, ruleID *int) (*User, error) {
	// Как и в CreateUser: UPDATE в CTE, чтобы JOIN был на свежих данных.
	query := `
	WITH upd AS (
		UPDATE users
		SET name = $1, female = $2, status_id = $3, rule_id = $4
		WHERE id = $5
		RETURNING id, name, female, status_id, rule_id
	)
	SELECT` + userColumns + ` FROM upd u` + userJoins

	u, err := scanUser(db.QueryRow(query, name, female, statusID, ruleID, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return u, nil
}
