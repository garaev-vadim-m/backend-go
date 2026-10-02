package rules

import "database/sql"

func GetAllRules(db *sql.DB) ([]*Rule, error) {
	query := `SELECT id,name,code FROM rules`

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rules []*Rule

	for rows.Next() {
		rule := Rule{}

		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Code,
		)

		if err != nil {
			return nil, err
		}

		rules = append(rules, &rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}
