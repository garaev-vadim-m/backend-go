package cart

import "database/sql"

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
