-- name: CreateCart :one
INSERT INTO carts (
    user_id
) VALUES (
    $1
) RETURNING *;

-- name: CreateGuestCart :one
INSERT INTO carts DEFAULT VALUES RETURNING *;

-- name: GetCartByID :one
SELECT * FROM carts WHERE id = $1;

-- name: GetCartByUserID :one
SELECT * FROM carts
WHERE user_id = $1 AND status = 'Active';

-- name: AddCartItem :one
INSERT INTO cart_items (
    cart_id, variant_id, qty, price, currency_code, thumbnail, product_name, product_path, variant_attributes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: IncrementCartItemQty :one
UPDATE cart_items
SET qty = qty + $3,
    updated_at = NOW()
WHERE cart_id = $1 AND variant_id = $2
RETURNING *;

-- name: UpdateCartItemQty :exec
UPDATE cart_items
SET qty = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteCartItem :exec
DELETE FROM cart_items
WHERE id = $1;

-- name: GetCartItemsByUserID :many
SELECT ci.* FROM cart_items ci
LEFT JOIN carts c ON ci.cart_id = c.id
WHERE  ci.is_removed IS FALSE AND c.user_id = $1;

-- name: GetCartItemsByCartID :many
SELECT * FROM cart_items
WHERE is_removed IS FALSE AND cart_id = $1;

-- name: UpdateCartClientSecret :exec
UPDATE carts
SET pi_client_secret = $3
WHERE id = $1 AND user_id = $2;

-- name: CountCartItems :one
SELECT COUNT(ci.*) FROM cart_items ci
LEFT JOIN carts c on c.id = ci.cart_id
WHERE user_id = $1 AND status = 'Active';

-- name: CountCartItemsByCartID :one
SELECT COUNT(*) FROM cart_items
WHERE cart_id = $1 AND is_removed = false;

-- name: CartItemExistsByUserID :one
SELECT EXISTS(
    SELECT 1 FROM cart_items ci
    LEFT JOIN carts c on ci.cart_id = c.id
    LEFT JOIN users u on u.id = c.user_id
    WHERE ci.id = $1 AND u.id = $2
);

-- name: CartItemExistsByCartID :one
SELECT EXISTS(
    SELECT 1 FROM cart_items
    WHERE id = $1 AND cart_id = $2
);