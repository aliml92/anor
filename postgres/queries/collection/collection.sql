-- name: GetCollectionByID :one
SELECT * FROM collections
WHERE id = $1;