package cart

import (
	"context"
	"database/sql"
)

func Get(ctx context.Context, db *sql.DB, userID int) (*Cart, error) {
	items, err := getItems(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	total, err := getTotal(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	return &Cart{Items: items, Total: total}, nil
}

func getItems(ctx context.Context, db *sql.DB, userID int) ([]Item, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT p.id, p.name, p.price::text, ci.quantity, (ci.quantity * p.price)::text
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id
		JOIN products p ON p.id = ci.product_id
		WHERE c.user_id = $1
		ORDER BY p.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}

	for rows.Next() {
		var item Item

		if err := rows.Scan(&item.ProductID, &item.Name, &item.Price, &item.Quantity, &item.Total); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func getTotal(ctx context.Context, db *sql.DB, userID int) (string, error) {
	var total string

	err := db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(ci.quantity * p.price), 0)::text
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id
		JOIN products p ON p.id = ci.product_id
		WHERE c.user_id = $1`, userID).Scan(&total)

	return total, err
}
