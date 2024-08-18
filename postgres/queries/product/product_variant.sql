-- name: CreateProductVariant :one
INSERT INTO product_variants (
    product_id, sku, qty, is_custom_priced, image_identifiers
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id;

-- name: GetProductVariantByID :one
SELECT * FROM product_variants
WHERE id = $1;

-- name: GetProductVariantsByProductID :many
SELECT
    pv.id AS variant_id,
    pv.product_id,
    pv.sku,
    pv.qty,
    pv.is_custom_priced,
    pv.image_identifiers,
    pv.created_at AS created_at,
    pv.updated_at AS updated_at,
    jsonb_object_agg(pa.attribute, pva.attribute_value) FILTER (WHERE pa.attribute IS NOT NULL) AS attributes
FROM
    product_variants pv
        LEFT JOIN
    product_variant_attributes pva ON pv.id = pva.variant_id
        LEFT JOIN
    product_attributes pa ON pva.product_attribute_id = pa.id
WHERE
    pv.product_id = $1
GROUP BY
    pv.id, pv.product_id, pv.sku, pv.qty, pv.is_custom_priced, pv.image_identifiers, pv.created_at, pv.updated_at;

-- name: GetProductVariantQtyInBulk :many
SELECT
    id,
    qty
FROM
    product_variants
WHERE
    id = ANY(sqlc.arg('variantIDs')::BIGINT[]);