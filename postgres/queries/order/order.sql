-- name: CreateOrder :one
INSERT INTO orders (
    cart_id, user_id, total_amount, payment_intent_id, shipping_address, billing_address
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: CreateOrderItems :copyfrom
INSERT INTO order_items (
   order_id, variant_id, qty, price, thumbnail, product_name, variant_attributes
) VALUES (
   $1, $2, $3, $4, $5, $6, $7
);

-- name: CountActiveOrders :one
SELECT COUNT(*) FROM orders
WHERE user_id = $1
  AND status IN ('Pending', 'Processing', 'Shipped');