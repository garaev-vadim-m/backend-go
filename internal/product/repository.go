package product

import "database/sql"

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
