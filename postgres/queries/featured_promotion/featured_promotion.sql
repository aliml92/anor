-- name: GetFeaturedPromotions :many
SELECT fp.*,
       CAST(
           CASE
               WHEN fp.type = 'collection' THEN c.name
               WHEN fp.type = 'category' THEN cat.category
               END AS TEXT
       ) AS target_name,
       CAST(
           CASE
               WHEN fp.type = 'collection' THEN c.handle
               WHEN fp.type = 'category' THEN cat.handle
               END AS TEXT
       ) AS target_handle
FROM featured_promotions fp
         LEFT JOIN collections c ON fp.type = 'collection' AND fp.target_id = c.id
         LEFT JOIN categories cat ON fp.type = 'category' AND fp.target_id = cat.id
WHERE (fp.end_date IS NULL OR fp.end_date >= CURRENT_DATE)
  AND (fp.start_date IS NULL OR fp.start_date <= CURRENT_DATE)
ORDER BY fp.display_order, fp.id
LIMIT $1;