package database

import (
	"database/sql"
	"fmt"
	"hello-go-backend/internal/models"
)

func GetAllStatus() ([]*models.Status, error) {
	query := `SELECT id, name, code FROM status`

	rows, err := DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var statuses []*models.Status

	for rows.Next() {
		status := models.Status{}

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

func GetStatus(id int) (*models.Status, error) {
	query := `SELECT id, name, code FROM status WHERE id = $1`

	status := &models.Status{}
	err := DB.QueryRow(query, id).Scan(
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

func CreateStatus(status *models.Status) (*models.Status, error) {
	query := `
	INSERT INTO status (name, code)
	VALUES ($1, $2)
	RETURNING id, name, code
	`

	err := DB.QueryRow(query, status.Name, status.Code).Scan(
		&status.ID,
		&status.Name,
		&status.Code,
	)

	if err != nil {
		return nil, err
	}

	return status, nil
}

func DeleteStatus(id int) error {
	query := `DELETE FROM status WHERE id = $1`

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

func UpdateStatus(id int, name string, code string) (*models.Status, error) {
	query := `
	UPDATE status
	SET name = $1, code = $2
	WHERE id = $3
	RETURNING id, name, code
	`

	status := &models.Status{}

	err := DB.QueryRow(query, name, code, id).Scan(
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
