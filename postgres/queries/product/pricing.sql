-- name: CreateProductPricing :exec
INSERT INTO product_pricing (
    product_id, base_price, currency_code, discount, discounted_price, is_on_sale
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: CreateProductVariantPricing :one
INSERT INTO product_variant_pricing (
    variant_id, base_price, currency_code, discount, discounted_price, is_on_sale
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING variant_id;

-- name: GetMinMaxPricesByCategoryID :one
SELECT
    MIN(pp.discounted_price)::numeric AS min_price,
    MAX(pp.discounted_price)::numeric AS max_price
FROM
    products p
        JOIN
    product_pricing pp ON p.id = pp.product_id
WHERE
    p.category_id = $1;

-- name: GetMinMaxPricesByCategoryIDs :one
SELECT
    MIN(pp.discounted_price)::numeric AS min_price,
    MAX(pp.discounted_price)::numeric AS max_price
FROM
    products p
        JOIN
    product_pricing pp ON p.id = pp.product_id
WHERE
    p.category_id = ANY(sqlc.arg('leafCategoryIds')::INT[]);

-- name: GetProductVariantPricingByVariantID :one
SELECT * FROM product_variant_pricing
WHERE variant_id = $1;

-- name: GetProductPricingByProductID :one
SELECT * FROM product_pricing
WHERE product_id = $1;