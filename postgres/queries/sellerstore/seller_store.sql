-- name: CreateStore :one
INSERT INTO seller_stores (
    public_id, seller_id, name, description
) VALUES (
    $1, $2, $3, $4
) RETURNING id;