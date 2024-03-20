-- name: CreateProduct :one
INSERT INTO products (
    store_id, category_id, name, brand, slug, short_info, image_links, specs, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id;

-- name: CreateProductPricing :exec
INSERT INTO product_pricing (
    product_id, base_price, currency_code, discount_level, discounted_amount, is_on_sale
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: CreateSKU :one
INSERT INTO skus (
    product_id, sku, quantity_in_stock, has_sep_pricing, image_refs
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id;		

-- name: CreateSKUPricing :one
INSERT INTO sku_pricing (
    sku_id, base_price, currency_code, discount_level, discounted_amount, is_on_sale
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING sku_id;	

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



