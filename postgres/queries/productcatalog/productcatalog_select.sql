-- name: GetAncestorCategoriesByCategoryID :many
WITH RECURSIVE categories_cte AS (
    SELECT c.id, c.category, c.slug, c.parent_id, 1 AS level
    FROM categories c
    WHERE c.id = $1
    UNION
    SELECT c.id, c.category, c.slug, c.parent_id, level + 1
    FROM categories c
    JOIN categories_cte cte ON cte.parent_id = c.id
)
SELECT *
FROM categories_cte
ORDER BY level DESC;

-- name: GetProductByID :one
SELECT sqlc.embed(p), sqlc.embed(s), sqlc.embed(pp) FROM products p
LEFT JOIN product_pricing pp ON pp.product_id = p.id
LEFT JOIN seller_stores s ON s.id = p.store_id
WHERE p.id = $1;

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


-- name: GetSiblingCategoriesByCategoryID :many
WITH category_cte AS (
    SELECT parent_id FROM categories c WHERE c.id = $1
)
SELECT * FROM categories WHERE parent_id = (SELECT parent_id FROM category_cte);

-- name: GetChildCategoriesByParentID :many
SELECT * FROM categories where parent_id = $1;

-- name: GetProductsByCategoryID :many
SELECT * FROM products
WHERE category_id = $1
OFFSET $2
LIMIT $3;

-- name: GetProductsByCategoryIDs :many
SELECT * FROM products
WHERE category_id = ANY(sqlc.arg('leafCategoryIds')::INT[])
OFFSET $1
LIMIT $2;

-- name: GetLeafCategoryIDsByCategoryID :many
WITH RECURSIVE categories_cte AS (
    SELECT c.id, c.category, c.slug, c.parent_id
    FROM categories c
    WHERE c.id = $1
    UNION ALL
    SELECT c2.id, c2.category, c2.slug, c2.parent_id
    FROM categories c2
             JOIN categories_cte cte ON cte.id = c2.parent_id
)
SELECT array_agg(id)::INT[] AS leaf_category_ids
FROM categories_cte
WHERE id NOT IN (SELECT parent_id FROM categories WHERE parent_id IS NOT NULL);

-- name: GetCategory :one
SELECT
    c.id,
    c.category,
    c.slug,
    c.parent_id,
    CASE
        WHEN EXISTS (SELECT 1 FROM categories WHERE parent_id = c.id) THEN false
        ELSE true
        END AS is_leaf_category
FROM
    categories c
WHERE
    c.id = $1;

