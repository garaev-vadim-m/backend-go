package database

import (
	"database/sql"
	"fmt"
	"hello-go-backend/internal/models"
)

func GetAllStatus() ([]*models.Status, error) {
	query := `SELECT id, name FROM status`

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
	qeury := `SELECT id, name FROM status WHERE id = $1`

	status := &models.Status{}
	err := DB.QueryRow(qeury, id).Scan(
		&status.ID,
		&status.Name,
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
	INSERT INTO status (name)
	VALUES ($1)
	RETURNING id, name
	`

	createStatus := &models.Status{}

	err := DB.QueryRow(query, status.Name).Scan(
		&status.ID,
		&status.Name,
	)

	if err != nil {
		return nil, err
	}

	return createStatus, nil
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
