-- name: CreateProduct :one
INSERT INTO products (
    store_id, category_id, name, brand, slug, short_info, image_urls, specs, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetProductByID :one
SELECT
    sqlc.embed(p),
    sqlc.embed(s),
    sqlc.embed(pp)
FROM
    products p
        LEFT JOIN
            product_pricing pp ON pp.product_id = p.id
        LEFT JOIN
            seller_stores s ON s.id = p.store_id
WHERE
    p.id = $1;


-- name: GetProductsByCategoryIDs :many
SELECT * FROM products
WHERE category_id = ANY(sqlc.arg('leafCategoryIds')::INT[])
OFFSET $1
LIMIT $2;