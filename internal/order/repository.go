package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hello-go-backend/internal/errs"
	"time"

	"github.com/lib/pq"
)

const uniqueKeyConstraint = "orders_user_id_idempotency_key_key"

type Order struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Total     string    `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}

type line struct {
	productID int
	qty       int
}

type querier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Checkout оформляет заказ из корзины. created=false значит, что заказ
// с этим ключом уже существовал, и вернулся он.

func Checkout(ctx context.Context, db *sql.DB, userID int, key string) (*Order, bool, error) {
	o, created, err := checkoutTx(ctx, db, userID, key)

	if err == nil {
		return o, created, nil
	}

	var pqErr *pq.Error

	if errors.As(err, &pqErr) && pqErr.Code == "23505" && pqErr.Constraint == uniqueKeyConstraint {
		existing, findErr := findByKey(ctx, db, userID, key)

		if findErr != nil {
			return nil, false, findErr
		}

		if existing != nil {
			return existing, false, nil
		}
	}

	return nil, false, err
}

func checkoutTx(ctx context.Context, db *sql.DB, userID int, key string) (*Order, bool, error) {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return nil, false, err
	}

	defer tx.Rollback()

	existing, err := findByKey(ctx, tx, userID, key)

	if err != nil {
		return nil, false, err
	}

	if existing != nil {
		return existing, false, nil
	}

	lines, err := lockCartLines(ctx, tx, userID)

	if err != nil {
		return nil, false, errs.ErrEmptyCart
	}

	if len(lines) == 0 {
		return nil, false, errs.ErrEmptyCart
	}

	for _, l := range lines {
		if err := reserve(ctx, tx, l.productID, l.qty); err != nil {
			return nil, false, err
		}
	}

	o := &Order{}

	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, idempotency_key) VALUES ($1, $2)
		RETURNING id, status, total::text, created_at`,
		userID, key).Scan(&o.ID, &o.Status, &o.Total, &o.CreatedAt)

	if err != nil {
		return nil, false, err
	}

	// Цена копируется из products в момент заказа (снимок).
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		SELECT $1::int, ci.product_id, ci.quantity, p.price
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id
		JOIN products p ON p.id = ci.product_id
		WHERE c.user_id = $2`, o.ID, userID); err != nil {
		return nil, false, err
	}

	// Итог считает БД: float64 для денег не используем.
	if err := tx.QueryRowContext(ctx, `
		UPDATE orders
		SET total = (SELECT COALESCE(SUM(quantity * price), 0) FROM order_items WHERE order_id = $1)
		WHERE id = $1
		RETURNING total::text`, o.ID).Scan(&o.Total); err != nil {
		return nil, false, err
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM cart_items
		WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1)`, userID); err != nil {
		return nil, false, err
	}

	if err := tx.Commit(); err != nil {
		return nil, false, err
	}
	return o, true, nil
}

func findByKey(ctx context.Context, q querier, userID int, key string) (*Order, error) {
	o := &Order{}

	err := q.QueryRowContext(ctx, `
		SELECT id, status, total::text, created_at
		FROM orders WHERE user_id = $1 AND idempotency_key = $2`,
		userID, key).Scan(&o.ID, &o.Status, &o.Total, &o.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return o, nil
}

// lockCartLines блокирует строки корзины до конца транзакции.
// ORDER BY фиксирует порядок блокировок и защищает от взаимных дедлоков.
func lockCartLines(ctx context.Context, tx *sql.Tx, userID int) ([]line, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT ci.product_id, ci.quantity
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id
		WHERE c.user_id = $1
		ORDER BY ci.product_id
		FOR UPDATE OF ci`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var lines []line

	for rows.Next() {
		var l line

		if err := rows.Scan(&l.productID, &l.qty); err != nil {
			return nil, err
		}

		lines = append(lines, l)
	}
	return lines, rows.Err()
}

// reserve списывает остаток одной строкой. Условие stock >= qty не даёт уйти в минус.
func reserve(ctx context.Context, tx *sql.Tx, productID, qty int) error {
	res, err := tx.ExecContext(ctx,
		`UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`,
		qty, productID)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("product %d: %w", productID, errs.ErrOutOfStock)
	}

	return nil
}
