package rules

import (
	"database/sql"
	"fmt"
)

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

func GetRule(db *sql.DB, id int) (*Rule, error) {
	query := `SELECT id, name, code FROM rules WHERE id = $1`

	rule := &Rule{}

	err := db.QueryRow(query, id).Scan(
		&rule.ID,
		&rule.Name,
		&rule.Code,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("status with id %d not found", id)
		}
		return nil, err
	}

	return rule, nil
}

func CreateRule(db *sql.DB, name string, code string) (*Rule, error) {
	query := `
	INSERT INTO rules (name, code)
	VALUES ($1, $2)
	RETURNING id, name, code
	`

	st := &Rule{}

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
