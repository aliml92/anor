-- name: IncrementCartItemQty :one
UPDATE cart_items
SET qty = CASE
              WHEN is_removed = TRUE THEN $3
              ELSE qty + $3
    END,
    is_removed = FALSE,
    updated_at = NOW()
WHERE cart_id = $1
  AND variant_id = $2
RETURNING *;

-- name: UpdateCartItemQty :execrows
UPDATE cart_items
SET qty = $2,
    updated_at = NOW()
WHERE id = $1 AND is_removed IS FALSE;

-- name: DeleteCartItem :execrows
UPDATE cart_items
SET is_removed = TRUE
WHERE id = $1;


-- name: UpdateCartStatus :execrows
UPDATE carts
SET status = $2
WHERE id = $1;
