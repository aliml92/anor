-- name: Create :one
INSERT INTO orders (
   user_id, cart_id, payment_method, payment_status, status, shipping_address_id, is_pickup, amount, currency
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: CreateOrderItems :copyfrom
INSERT INTO order_items (
    order_id, variant_id, qty, price, thumbnail, product_name, variant_attributes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetOrders :many
SELECT * FROM orders
WHERE cart_id = $1;

-- name: Get :one
SELECT * FROM orders
WHERE id = $1;

-- name: ListActiveByUserID :many
SELECT *
FROM orders
WHERE user_id = $1
  AND (
    (payment_status = 'Paid' AND status IN ('Processing', 'Shipped', 'Delivered'))
        OR
    (payment_status = 'Pending' AND payment_method = 'PayOnDelivery')
    )
ORDER BY created_at DESC
LIMIT $2
    OFFSET $3;

-- name: ListActiveWithPaymentAndAddressByUserID :many
SELECT
    sqlc.embed(o),
    sqlc.embed(a),
    sqlc.embed(scp)
FROM
    orders o
        LEFT JOIN
    addresses a ON o.shipping_address_id = a.id
        LEFT JOIN
    stripe_card_payments scp ON o.id = scp.order_id
WHERE
    o.user_id = $1
  AND (
    (o.payment_status = 'Paid' AND o.status IN ('Processing', 'Shipped', 'Delivered'))
        OR
    (o.payment_status = 'Pending' AND o.payment_method = 'PayOnDelivery')
    )
ORDER BY
    o.created_at DESC
LIMIT $2
    OFFSET $3;

-- name: ListUnpaidPaymentByUserID :many
SELECT *
FROM orders
WHERE user_id = $1
  AND payment_status = 'Pending'
  AND payment_method IN ('Stripe', 'Paypal', 'AnorInstallment')
ORDER BY created_at DESC
LIMIT $2
    OFFSET $3;

-- name: ListUnpaidWithPaymentAndAddressByUserID :many
SELECT
    sqlc.embed(o),
    sqlc.embed(a),
    sqlc.embed(scp)
FROM
    orders o
        LEFT JOIN
    addresses a ON o.shipping_address_id = a.id
        LEFT JOIN
    stripe_card_payments scp ON o.id = scp.order_id
WHERE
    o.user_id = $1
  AND payment_status = 'Pending'
  AND payment_method IN ('Stripe', 'Paypal', 'AnorInstallment')
ORDER BY
    o.created_at DESC
LIMIT $2
    OFFSET $3;

-- name: ListByUserID :many
SELECT *
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
    OFFSET $3;

-- name: ListWithPaymentAndAddressByUserID :many
SELECT
    sqlc.embed(o),
    sqlc.embed(a),
    sqlc.embed(scp)
FROM
    orders o
        LEFT JOIN
    addresses a ON o.shipping_address_id = a.id
        LEFT JOIN
    stripe_card_payments scp ON o.id = scp.order_id
WHERE o.user_id = $1
ORDER BY o.created_at DESC
LIMIT $2
    OFFSET $3;

-- name: GetItemByOrderIds :many
SELECT * FROM order_items
WHERE order_id = ANY(sqlc.arg('orderIDs')::BIGINT[]);