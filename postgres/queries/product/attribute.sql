-- name: CreateProductAttribute :one
INSERT INTO product_attributes (
    product_id, attribute
) VALUES (
    $1, $2
) RETURNING *;

-- name: CreateProductVariantAttributeValues :exec
INSERT INTO product_variant_attributes (
    variant_id, product_attribute_id, attribute_value
) VALUES (
    $1, $2, $3
);

-- name: GetProductAttributesByProductID :many
SELECT
    pa.attribute AS attribute_name,
    ARRAY_AGG(DISTINCT pva.attribute_value)::text[] AS attribute_values
FROM
    product_attributes pa
        JOIN
            product_variant_attributes pva ON pa.id = pva.product_attribute_id
WHERE
    pa.product_id = $1
GROUP BY
    pa.attribute;

-- name: GetProductAttributesByCategoryID :many
SELECT DISTINCT
    pa.attribute,
    array_agg(DISTINCT pva.attribute_value ORDER BY pva.attribute_value)::text[] AS values
FROM
    products p
        JOIN
    product_attributes pa ON p.id = pa.product_id
        JOIN
    product_variants pv ON p.id = pv.product_id
        JOIN
    product_variant_attributes pva ON pv.id = pva.variant_id
WHERE
    p.category_id = $1
GROUP BY
    pa.attribute;

-- name: GetProductAttributesByCategoryIDs :many
SELECT DISTINCT
    pa.attribute,
    array_agg(DISTINCT pva.attribute_value ORDER BY pva.attribute_value)::text[] AS values
FROM
    products p
        JOIN
    product_attributes pa ON p.id = pa.product_id
        JOIN
    product_variants pv ON p.id = pv.product_id
        JOIN
    product_variant_attributes pva ON pv.id = pva.variant_id
WHERE
    p.category_id = ANY(sqlc.arg('leafCategoryIds')::INT[])
GROUP BY
    pa.attribute;