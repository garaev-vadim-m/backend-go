package status

import (
	"database/sql"
	"fmt"
)

func GetAllStatus(db *sql.DB) ([]*Status, error) {
	query := `SELECT id, name, code FROM status`

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var statuses []*Status

	for rows.Next() {
		status := Status{}

		err := rows.Scan(
			&status.ID,
			&status.Name,
			&status.Code,
		)

		if err != nil {
			return nil, err
		}

		statuses = append(statuses, &status)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statuses, nil
}

func GetStatus(db *sql.DB, id int) (*Status, error) {
	query := `SELECT id, name, code FROM status WHERE id = $1`

	status := &Status{}
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

func CreateStatus(db *sql.DB, name string, code string) (*Status, error) {
	query := `
	INSERT INTO status (name, code)
	VALUES ($1, $2)
	RETURNING id, name, code
	`

	st := &Status{}
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
	query := `DELETE FROM status WHERE id = $1`

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

func UpdateStatus(db *sql.DB, id int, name string, code string) (*Status, error) {
	query := `
	UPDATE status
	SET name = $1, code = $2
	WHERE id = $3
	RETURNING id, name, code
	`

	status := &Status{}

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
