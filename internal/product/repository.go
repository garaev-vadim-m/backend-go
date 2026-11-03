package product

import (
	"database/sql"
	"fmt"

	"hello-go-backend/internal/status"
)

// Общие части SELECT: колонки и JOIN в одном месте.
// Алиасы везде одинаковые: p (products) и st (status).
const productColumns = `
	p.id, p.name, p.description,
	st.id, st.name, st.code`

const productJoins = `
	LEFT JOIN status st ON st.id = p.status_id`

type scanner interface {
	Scan(dest ...any) error
}

// scanProduct разбирает строку (с JOIN-нутым статусом) в Product.
// Если статуса нет, status_id и status остаются nil.
func scanProduct(s scanner) (*Product, error) {
	p := &Product{}

	var statusID sql.NullInt64
	var statusName, statusCode sql.NullString

	err := s.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&statusID,
		&statusName,
		&statusCode,
	)
	if err != nil {
		return nil, err
	}

	if statusID.Valid {
		id := int(statusID.Int64)
		p.Status = &status.Status{
			ID:   id,
			Name: statusName.String,
			Code: statusCode.String,
		}
	}

	return p, nil
}

func GetAllProduct(db *sql.DB) ([]*Product, error) {
	query := `SELECT` + productColumns + ` FROM products p` + productJoins

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product

	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func GetProduct(db *sql.DB, id int) (*Product, error) {
	query := `SELECT` + productColumns + ` FROM products p` + productJoins + ` WHERE p.id = $1`

	p, err := scanProduct(db.QueryRow(query, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, err
	}

	return p, nil
}

func CreateProduct(db *sql.DB, name string, description *string, statusId *int) (*Product, error) {
	// RETURNING не умеет JOIN, поэтому вставляем в CTE и джойним его.
	query := `
	WITH ins AS (
		INSERT INTO products (name, description, status_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, status_id
	)
	SELECT` + productColumns + ` FROM ins p` + productJoins

	p, err := scanProduct(db.QueryRow(query, name, description, statusId))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func UpdateProduct(db *sql.DB, id int, name string, description *string, statusId *int) (*Product, error) {
	query := `
	WITH upd AS (
		UPDATE products
		SET name = $1, description = $2, status_id = $3
		WHERE id = $4
		RETURNING id, name, description, status_id
	)
	SELECT` + productColumns + ` FROM upd p` + productJoins

	p, err := scanProduct(db.QueryRow(query, name, description, statusId, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, err
	}

	return p, nil
}

func DeleteProduct(db *sql.DB, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	return nil
}
