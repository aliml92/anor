-- name: CreateProduct :one
INSERT INTO products (
    store_id, category_id, name, brand, handle, short_information, image_urls, specifications, status
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
            stores s ON s.id = p.store_id
WHERE
    p.id = $1;

-- name: GetProductBrandsByCategoryID :many
SELECT array_agg(DISTINCT p.brand ORDER BY p.brand)::text[] AS brands FROM products p
WHERE p.category_id = $1;

-- name: GetProductBrandsByCategoryIDs :many
SELECT array_agg(DISTINCT p.brand ORDER BY p.brand)::text[] AS brands FROM products p
WHERE p.category_id = ANY(sqlc.arg('leafCategoryIds')::INT[]);

-- name: GetProductImageURLsByID :one
SELECT image_urls
FROM products
WHERE id = $1;

-- name: CollectCartData :one
SELECT
    p.id AS product_id,
    p.name AS product_name,
    p.handle AS product_handle,
    p.image_urls AS product_image_urls,
    pv.sku,
    pv.qty,
    pv.is_custom_priced,
    pv.image_identifiers,
    COALESCE(
                    json_object_agg(
                    pa.attribute,
                    pva.attribute_value
                                   ) FILTER (WHERE pa.id IS NOT NULL),
                    '{}'
    )::json AS variant_attributes
FROM
    product_variants pv
        JOIN
    products p ON pv.product_id = p.id
        LEFT JOIN
    product_variant_attributes pva ON pv.id = pva.variant_id
        LEFT JOIN
    product_attributes pa ON pva.product_attribute_id = pa.id
WHERE
    pv.id = $1
GROUP BY
    p.id, p.name, p.image_urls, pv.sku, pv.qty, pv.is_custom_priced, pv.image_identifiers;

-- name: GetProductsByCreatedAtDesc :many
SELECT sqlc.embed(p), sqlc.embed(pp) FROM products p
LEFT JOIN product_pricing pp ON p.id = pp.product_id
ORDER BY p.created_at DESC
LIMIT $1;