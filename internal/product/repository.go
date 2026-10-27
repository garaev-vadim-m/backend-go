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

func CreateProduct(db *sql.DB, name string, description *string, statusId *int) (*Product, error) {
	query := `
	INSERT INTO products (name, description, status_id)
	VALUES ($1, $2, $3)
	RETURNING id, name, description, status_id
	`

	p := &Product{}
	err := db.QueryRow(query, name, description, statusId).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.StatusId,
	)

	if err != nil {
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
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func UpdateProduct(db *sql.DB, id int, name string, description *string, statusId *int) (*Product, error) {
	query := `
	UPDATE products
	SET name = $1, description = $2, status_id = $3
	WHERE id = $4
	RETURNING id, name, description, status_id
	`

	product := &Product{}

	err := db.QueryRow(query, name, description, statusId, id).Scan(
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
