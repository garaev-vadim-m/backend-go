package userstatus

import (
	"database/sql"
	"fmt"
)

func GetAllStatus(db *sql.DB) ([]*UserStatus, error) {
	query := `SELECT id, name, code FROM user_status`

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user_statuses []*UserStatus

	for rows.Next() {
		user_status := UserStatus{}

		err := rows.Scan(
			&user_status.ID,
			&user_status.Name,
			&user_status.Code,
		)

		if err != nil {
			return nil, err
		}

		user_statuses = append(user_statuses, &user_status)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return user_statuses, nil
}

func GetStatus(db *sql.DB, id int) (*UserStatus, error) {
	query := `SELECT id, name, code FROM user_status WHERE id = $1`

	status := &UserStatus{}
	err := db.QueryRow(query, id).Scan(
		&status.ID,
		&status.Name,
		&status.Code,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("status with id %d not found", id)
		}
		return nil, err
	}

	return status, nil
}

func CreateStatus(db *sql.DB, name string, code string) (*UserStatus, error) {
	query := `
	INSERT INTO user_status (name, code)
	VALUES ($1, $2)
	RETURNING id, name, code
	`

	st := &UserStatus{}
	err := db.QueryRow(query, name, code).Scan(
		&st.ID,
		&st.Name,
		&st.Code,
	)

	if err != nil {
		return nil, err
	}

	return st, nil
}

func DeleteStatus(db *sql.DB, id int) error {
	query := `DELETE FROM user_status WHERE id = $1`

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

func UpdateStatus(db *sql.DB, id int, name string, code string) (*UserStatus, error) {
	query := `
	UPDATE status
	SET name = $1, code = $2
	WHERE id = $3
	RETURNING id, name, code
	`

	status := &UserStatus{}

	err := db.QueryRow(query, name, code, id).Scan(
		&status.ID,
		&status.Name,
		&status.Code,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("status with id %d not found", id)
		}
		return nil, err
	}

	return status, nil
}
