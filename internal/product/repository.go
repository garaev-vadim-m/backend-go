package product

import (
	"database/sql"
	"fmt"
)

func GetAllProduct(db *sql.DB) ([]*Product, error) {
	query := `SELECT id, name, description, status_id FROM products`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*Product

	for rows.Next() {
		product := Product{}

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.StatusId,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func GetProduct(db *sql.DB, id int) (*Product, error) {
	query := `SELECT id, name, description, status_id FROM products WHERE id = $1`

	product := &Product{}

	err := db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.StatusId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, err
	}

	return product, nil
}
