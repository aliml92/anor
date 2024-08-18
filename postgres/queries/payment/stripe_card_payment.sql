-- name: GetByOrderID :one
SELECT * FROM stripe_card_payments
WHERE order_id = $1;

-- name: Create :exec
INSERT INTO stripe_card_payments
    (order_id, user_id, billing_address_id, payment_intent_id, payment_method_id, amount, currency, status, client_secret, last_error, card_last4, card_brand, created_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);


