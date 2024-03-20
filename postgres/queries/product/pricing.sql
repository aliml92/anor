-- name: CreateProductPricing :exec
INSERT INTO product_pricing (
    product_id, base_price, currency_code, discount_level, discounted_amount, is_on_sale
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: CreateSKUPricing :one
INSERT INTO sku_pricing (
    sku_id, base_price, currency_code, discount_level, discounted_amount, is_on_sale
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING sku_id;