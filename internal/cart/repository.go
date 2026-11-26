package cart

import (
	"context"
	"database/sql"
)

// AddItem добавляет товар в корзину, корзину создает при необходимости.
// Повторное добавление увеличивает количество.
func AddItem(db *sql.DB, userID, productID, qty int) error {
	_, err := db.Exec(`
		WITH c AS (
			INSERT INTO carts (user_id) VALUES ($1)
			ON CONFLICT (user_id) DO UPDATE SET updated_at = now()
			RETURNING id
		)
		INSERT INTO cart_items (cart_id, product_id, quantity)
		SELECT id, $2::int, $3::int FROM c
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`,
		userID, productID, qty)
	return err
}

// GetCart возвращает корзину одним запросом: позиции и итог из одного снимка.
// Пустая корзина: items = [], total = "0.00".
func GetCart(ctx context.Context, db *sql.DB, userID int) (*Cart, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT p.id, p.name, p.price::text, ci.quantity,
		       (ci.quantity * p.price)::text,
		       (SUM(ci.quantity * p.price) OVER ())::text
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id
		JOIN products p ON p.id = ci.product_id
		WHERE c.user_id = $1
		ORDER BY p.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart := &Cart{Items: []Item{}, Total: "0.00"}

	for rows.Next() {
		var item Item
		var total string

		if err := rows.Scan(
			&item.ProductID,
			&item.Name,
			&item.Price,
			&item.Quantity,
			&item.Total,
			&total,
		); err != nil {
			return nil, err
		}

		cart.Items = append(cart.Items, item)
		cart.Total = total // одинаковое значение во всех строках, берём последнее
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveItem удаляет товар из корзины пользователя.
// Возвращает false, если такого товара в корзине не было.
func RemoveItem(db *sql.DB, userID, productID int) (bool, error) {
	res, err := db.Exec(`
		DELETE FROM cart_items ci
		USING carts c
		WHERE ci.cart_id = c.id
		  AND c.user_id = $1
		  AND ci.product_id = $2`, userID, productID)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
