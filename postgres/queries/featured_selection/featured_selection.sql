-- name: Create :one
INSERT INTO featured_selections (
    resource_path, banner_info, image_url, query_params, start_date, end_date, display_order
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListAllActive :many
SELECT * FROM featured_selections
WHERE (start_date IS NULL OR start_date <= CURRENT_DATE)
  AND (end_date IS NULL OR end_date >= CURRENT_DATE)
ORDER BY display_order ASC;
