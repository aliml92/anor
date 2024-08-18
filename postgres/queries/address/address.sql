-- name: CreateAddress :one
INSERT INTO addresses (
    user_id, default_for, name, address_line_1, address_line_2, city, state_province, postal_code, country
) VALUES  (
  $1, $2, $3, $4, $5,$6,$7,$8, $9
) RETURNING *;

-- name: GetAddressByID :one
SELECT * FROM addresses
WHERE id = $1;

-- name: GetUserDefaultAddress :one
SELECT * FROM addresses
WHERE user_id = $1
    AND default_for = $2
ORDER BY created_at
OFFSET 1;

-- name: ListAddressesByUserID :many
SELECT * FROM addresses
WHERE user_id = $1
ORDER BY created_at
OFFSET $2
    LIMIT $3;

