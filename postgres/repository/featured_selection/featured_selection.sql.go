// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: featured_selection.sql

package featured_selection

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const create = `-- name: Create :one
INSERT INTO featured_selections (
    resource_path, banner_info, image_url, query_params, start_date, end_date, display_order
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, resource_path, banner_info, image_url, query_params, start_date, end_date, display_order, created_at, updated_at
`

func (q *Queries) Create(ctx context.Context, resourcePath string, bannerInfo []byte, imageUrl string, queryParams []byte, startDate pgtype.Date, endDate pgtype.Date, displayOrder *int32) (*FeaturedSelection, error) {
	row := q.db.QueryRow(ctx, create,
		resourcePath,
		bannerInfo,
		imageUrl,
		queryParams,
		startDate,
		endDate,
		displayOrder,
	)
	var i FeaturedSelection
	err := row.Scan(
		&i.ID,
		&i.ResourcePath,
		&i.BannerInfo,
		&i.ImageUrl,
		&i.QueryParams,
		&i.StartDate,
		&i.EndDate,
		&i.DisplayOrder,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const listAllActive = `-- name: ListAllActive :many
SELECT id, resource_path, banner_info, image_url, query_params, start_date, end_date, display_order, created_at, updated_at FROM featured_selections
WHERE (start_date IS NULL OR start_date <= CURRENT_DATE)
  AND (end_date IS NULL OR end_date >= CURRENT_DATE)
ORDER BY display_order ASC
`

func (q *Queries) ListAllActive(ctx context.Context) ([]*FeaturedSelection, error) {
	rows, err := q.db.Query(ctx, listAllActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*FeaturedSelection
	for rows.Next() {
		var i FeaturedSelection
		if err := rows.Scan(
			&i.ID,
			&i.ResourcePath,
			&i.BannerInfo,
			&i.ImageUrl,
			&i.QueryParams,
			&i.StartDate,
			&i.EndDate,
			&i.DisplayOrder,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
