// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: order.sql

package order

import (
	"context"

	"github.com/shopspring/decimal"
)

const countActiveOrders = `-- name: CountActiveOrders :one
SELECT COUNT(*) FROM orders
WHERE user_id = $1
  AND status IN ('Pending', 'Processing', 'Shipped')
`

func (q *Queries) CountActiveOrders(ctx context.Context, userID *int64) (int64, error) {
	row := q.db.QueryRow(ctx, countActiveOrders, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (
    cart_id, user_id, total_amount, payment_intent_id, shipping_address, billing_address
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, cart_id, user_id, status, total_amount, payment_intent_id, shipping_address, billing_address, created_at, updated_at
`

func (q *Queries) CreateOrder(ctx context.Context, cartID int64, userID *int64, totalAmount decimal.Decimal, paymentIntentID string, shippingAddress []byte, billingAddress []byte) (*Order, error) {
	row := q.db.QueryRow(ctx, createOrder,
		cartID,
		userID,
		totalAmount,
		paymentIntentID,
		shippingAddress,
		billingAddress,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.CartID,
		&i.UserID,
		&i.Status,
		&i.TotalAmount,
		&i.PaymentIntentID,
		&i.ShippingAddress,
		&i.BillingAddress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
