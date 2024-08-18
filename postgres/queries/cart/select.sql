-- name: GetCartByID :one
SELECT * FROM carts
WHERE id = $1
  AND status = 'Open'
ORDER BY updated_at
LIMIT 1;

-- name: GetCartByUserID :one
SELECT * FROM carts
WHERE user_id = $1
  AND status = 'Open'
ORDER BY updated_at
LIMIT 1;


-- name: GetCartItemsByUserID :many
SELECT ci.* FROM cart_items ci
LEFT JOIN carts c ON ci.cart_id = c.id
WHERE c.user_id = $1
  AND ci.is_removed IS FALSE;

-- name: GetCartItemsByCartID :many
SELECT * FROM cart_items
WHERE cart_id = $1 AND is_removed IS FALSE;

-- name: CountCartItemsByUserIdAndCartStatus :one
SELECT COUNT(ci.*) FROM cart_items ci
                            LEFT JOIN carts c on c.id = ci.cart_id
WHERE user_id = $1
  AND status = $2
  AND ci.is_removed IS FALSE;

-- name: CountCartItemsByCartID :one
SELECT COUNT(*) FROM cart_items
WHERE cart_id = $1 AND is_removed = false;

-- name: IsCartItemExistsByCartUserID :one
SELECT EXISTS(
    SELECT 1 FROM cart_items ci
      LEFT JOIN carts c on ci.cart_id = c.id
    WHERE ci.id = $1
      AND c.user_id = $2
      AND c.status = 'Open'
);

-- name: CartItemExistsByCartID :one
SELECT EXISTS(
    SELECT 1 FROM cart_items
    WHERE id = $1 AND cart_id = $2
);