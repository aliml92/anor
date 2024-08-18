-- name: CreateCart :one
INSERT INTO carts (
    user_id
) VALUES (
    $1
) RETURNING *;

-- name: CreateGuestCart :one
INSERT INTO carts (
    expires_at
) VALUES (
    $1
) RETURNING *;

-- name: AddCartItem :one
INSERT INTO cart_items (
    cart_id, variant_id, qty, price, currency, thumbnail, product_name, product_path, variant_attributes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: MergeGuestCartWithUserCart :exec
INSERT INTO cart_items (cart_id, variant_id, qty, price, thumbnail, product_name, product_path, variant_attributes, is_removed, created_at, updated_at)
SELECT
    user_cart.id, -- Use the user's cart id
    guest_items.variant_id,
    CASE
        WHEN user_items.qty IS NULL THEN guest_items.qty
        WHEN user_items.qty >= guest_items.qty THEN user_items.qty
        ELSE guest_items.qty
        END AS qty,
    guest_items.price,
    guest_items.thumbnail,
    guest_items.product_name,
    guest_items.product_path,
    guest_items.variant_attributes,
    guest_items.is_removed,
    LEAST(guest_items.created_at, COALESCE(user_items.created_at, guest_items.created_at)) AS created_at,
    CURRENT_TIMESTAMP AS updated_at
FROM
    cart_items guest_items
        JOIN
    carts guest_cart ON guest_items.cart_id = guest_cart.id
        CROSS JOIN
    carts user_cart
        LEFT JOIN
    cart_items user_items ON user_items.cart_id = user_cart.id AND user_items.variant_id = guest_items.variant_id
WHERE
    guest_cart.id = sqlc.arg('guestCartID')  -- Use $1 for guest_cart_id
  AND user_cart.id = sqlc.arg('userCartID')  -- Use $2 for user_cart_id
  AND (user_items.id IS NULL OR guest_items.qty > user_items.qty)
ON CONFLICT (cart_id, variant_id)
    DO UPDATE SET
                  qty = EXCLUDED.qty,
                  updated_at = CURRENT_TIMESTAMP;