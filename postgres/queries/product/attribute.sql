-- name: CreateProductAttribute :one
INSERT INTO product_attributes (
    product_id, attribute
) VALUES (
    $1, $2
) RETURNING id;

-- name: CreateSKUProductAttributeValues :exec
INSERT INTO sku_product_attribute_values (
    sku_id, product_attribute_id, attribute_value
) VALUES (
    $1, $2, $3
);

-- name: GetProductAttributesByProductID :many
SELECT
    pa.attribute AS attribute_name,
    ARRAY_AGG(DISTINCT spav.attribute_value)::text[] AS attribute_values
FROM
    product_attributes pa
        JOIN
            sku_product_attribute_values spav ON pa.id = spav.product_attribute_id
WHERE
    pa.product_id = $1
GROUP BY
    pa.attribute;