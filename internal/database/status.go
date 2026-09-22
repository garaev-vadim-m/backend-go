package database

import "hello-go-backend/internal/models"

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
