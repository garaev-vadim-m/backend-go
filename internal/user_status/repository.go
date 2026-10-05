package userstatus

import "database/sql"

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
