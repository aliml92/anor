-- name: CreateStore :one
INSERT INTO stores (
    handle, user_id, name, description
) VALUES (
    $1, $2, $3, $4
) RETURNING id;