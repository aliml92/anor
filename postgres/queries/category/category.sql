-- name: CreateTopCategory :one
INSERT INTO categories (
    category, handle
) VALUES (
    $1, $2
) RETURNING *;


-- name: CreateChildCategoryIfNotExists :one
WITH inserted AS (
    INSERT INTO categories
        (category, handle, parent_id)
    VALUES
        ($1, $2, $3)
    ON CONFLICT (category, parent_id)
        DO NOTHING
    RETURNING *
)
SELECT * FROM inserted
UNION ALL
SELECT * FROM categories
WHERE category = $1 AND parent_id = $3;


-- name: GetCategoryWithAncestors :many
WITH RECURSIVE categories_cte AS (
    SELECT c.id, c.category, c.handle, c.parent_id, 1 AS level
    FROM categories c
    WHERE c.id = $1
    UNION
    SELECT c.id, c.category, c.handle, c.parent_id, level + 1
    FROM categories c
             JOIN categories_cte cte ON cte.parent_id = c.id
)
SELECT *
FROM categories_cte
ORDER BY level DESC;


-- name: GetCategoryWithSiblings :many
WITH category_cte AS (
    SELECT parent_id FROM categories c WHERE c.id = $1
)
SELECT * FROM categories WHERE parent_id = (SELECT parent_id FROM category_cte);


-- name: GetChildCategories :many
SELECT * FROM categories where parent_id = $1;


-- name: GetLeafCategoryIDs :many
WITH RECURSIVE categories_cte AS (
    SELECT c.id, c.category, c.handle, c.parent_id
    FROM categories c
    WHERE c.id = $1
    UNION ALL
    SELECT c2.id, c2.category, c2.handle, c2.parent_id
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
    c.handle,
    c.parent_id,
    CASE
        WHEN EXISTS (SELECT 1 FROM categories WHERE parent_id = c.id) THEN false
        ELSE true
        END AS is_leaf_category
FROM
    categories c
WHERE
    c.id = $1;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE category = $1
ORDER BY id
LIMIT 1;

-- name: GetRootCategories :many
SELECT * FROM categories WHERE parent_id IS NULL;