// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: select.sql

package cart

import (
	"context"
)

const cartItemExistsByCartID = `-- name: CartItemExistsByCartID :one
SELECT EXISTS(
    SELECT 1 FROM cart_items
    WHERE id = $1 AND cart_id = $2
)
`

func (q *Queries) CartItemExistsByCartID(ctx context.Context, iD int64, cartID int64) (bool, error) {
	row := q.db.QueryRow(ctx, cartItemExistsByCartID, iD, cartID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const countCartItemsByCartID = `-- name: CountCartItemsByCartID :one
SELECT COUNT(*) FROM cart_items
WHERE cart_id = $1 AND is_removed = false
`

func (q *Queries) CountCartItemsByCartID(ctx context.Context, cartID int64) (int64, error) {
	row := q.db.QueryRow(ctx, countCartItemsByCartID, cartID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countCartItemsByUserIdAndCartStatus = `-- name: CountCartItemsByUserIdAndCartStatus :one
SELECT COUNT(ci.*) FROM cart_items ci
                            LEFT JOIN carts c on c.id = ci.cart_id
WHERE user_id = $1
  AND status = $2
  AND ci.is_removed IS FALSE
`

func (q *Queries) CountCartItemsByUserIdAndCartStatus(ctx context.Context, userID *int64, status CartStatus) (int64, error) {
	row := q.db.QueryRow(ctx, countCartItemsByUserIdAndCartStatus, userID, status)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getCartByID = `-- name: GetCartByID :one
SELECT id, user_id, status, last_activity, expires_at, created_at, updated_at FROM carts
WHERE id = $1
  AND status = 'Open'
ORDER BY updated_at
LIMIT 1
`

func (q *Queries) GetCartByID(ctx context.Context, id int64) (*Cart, error) {
	row := q.db.QueryRow(ctx, getCartByID, id)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.LastActivity,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getCartByUserID = `-- name: GetCartByUserID :one
SELECT id, user_id, status, last_activity, expires_at, created_at, updated_at FROM carts
WHERE user_id = $1
  AND status = 'Open'
ORDER BY updated_at
LIMIT 1
`

func (q *Queries) GetCartByUserID(ctx context.Context, userID *int64) (*Cart, error) {
	row := q.db.QueryRow(ctx, getCartByUserID, userID)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.LastActivity,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getCartItemsByCartID = `-- name: GetCartItemsByCartID :many
SELECT id, cart_id, variant_id, qty, price, currency, thumbnail, product_name, product_path, variant_attributes, is_removed, created_at, updated_at FROM cart_items
WHERE cart_id = $1 AND is_removed IS FALSE
`

func (q *Queries) GetCartItemsByCartID(ctx context.Context, cartID int64) ([]*CartItem, error) {
	rows, err := q.db.Query(ctx, getCartItemsByCartID, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*CartItem
	for rows.Next() {
		var i CartItem
		if err := rows.Scan(
			&i.ID,
			&i.CartID,
			&i.VariantID,
			&i.Qty,
			&i.Price,
			&i.Currency,
			&i.Thumbnail,
			&i.ProductName,
			&i.ProductPath,
			&i.VariantAttributes,
			&i.IsRemoved,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCartItemsByUserID = `-- name: GetCartItemsByUserID :many
SELECT ci.id, ci.cart_id, ci.variant_id, ci.qty, ci.price, ci.currency, ci.thumbnail, ci.product_name, ci.product_path, ci.variant_attributes, ci.is_removed, ci.created_at, ci.updated_at FROM cart_items ci
LEFT JOIN carts c ON ci.cart_id = c.id
WHERE c.user_id = $1
  AND ci.is_removed IS FALSE
`

func (q *Queries) GetCartItemsByUserID(ctx context.Context, userID *int64) ([]*CartItem, error) {
	rows, err := q.db.Query(ctx, getCartItemsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*CartItem
	for rows.Next() {
		var i CartItem
		if err := rows.Scan(
			&i.ID,
			&i.CartID,
			&i.VariantID,
			&i.Qty,
			&i.Price,
			&i.Currency,
			&i.Thumbnail,
			&i.ProductName,
			&i.ProductPath,
			&i.VariantAttributes,
			&i.IsRemoved,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const isCartItemExistsByCartUserID = `-- name: IsCartItemExistsByCartUserID :one
SELECT EXISTS(
    SELECT 1 FROM cart_items ci
      LEFT JOIN carts c on ci.cart_id = c.id
    WHERE ci.id = $1
      AND c.user_id = $2
      AND c.status = 'Open'
)
`

func (q *Queries) IsCartItemExistsByCartUserID(ctx context.Context, iD int64, userID *int64) (bool, error) {
	row := q.db.QueryRow(ctx, isCartItemExistsByCartUserID, iD, userID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
