-- name: CreateSKU :one
INSERT INTO skus (
    product_id, sku, quantity_in_stock, has_sep_pricing, image_refs
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id;

-- name: GetSKUsByProductID :many
SELECT
    s.id AS sku_id,
    s.product_id,
    s.sku,
    s.quantity_in_stock,
    s.has_sep_pricing,
    s.image_refs,
    s.created_at AS sku_created_at,
    s.updated_at AS sku_updated_at,
    jsonb_object_agg(pa.attribute, spav.attribute_value) FILTER (WHERE pa.attribute IS NOT NULL) AS attributes
FROM
    skus s
        LEFT JOIN
    sku_product_attribute_values spav ON s.id = spav.sku_id
        LEFT JOIN
    product_attributes pa ON spav.product_attribute_id = pa.id
WHERE
    s.product_id = $1
GROUP BY
    s.id, s.product_id, s.sku, s.quantity_in_stock, s.has_sep_pricing, s.image_refs, s.created_at, s.updated_at;