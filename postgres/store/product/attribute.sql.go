// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: attribute.sql

package product

import (
	"context"
)

const createProductAttribute = `-- name: CreateProductAttribute :one
INSERT INTO product_attributes (
    product_id, attribute
) VALUES (
    $1, $2
) RETURNING id
`

func (q *Queries) CreateProductAttribute(ctx context.Context, productID int64, attribute string) (int64, error) {
	row := q.db.QueryRow(ctx, createProductAttribute, productID, attribute)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createSKUProductAttributeValues = `-- name: CreateSKUProductAttributeValues :exec
INSERT INTO sku_product_attribute_values (
    sku_id, product_attribute_id, attribute_value
) VALUES (
    $1, $2, $3
)
`

func (q *Queries) CreateSKUProductAttributeValues(ctx context.Context, skuID int64, productAttributeID int64, attributeValue string) error {
	_, err := q.db.Exec(ctx, createSKUProductAttributeValues, skuID, productAttributeID, attributeValue)
	return err
}

const getProductAttributesByProductID = `-- name: GetProductAttributesByProductID :many
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
    pa.attribute
`

type GetProductAttributesByProductIDRow struct {
	AttributeName   string
	AttributeValues []string
}

func (q *Queries) GetProductAttributesByProductID(ctx context.Context, productID int64) ([]*GetProductAttributesByProductIDRow, error) {
	rows, err := q.db.Query(ctx, getProductAttributesByProductID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetProductAttributesByProductIDRow
	for rows.Next() {
		var i GetProductAttributesByProductIDRow
		if err := rows.Scan(&i.AttributeName, &i.AttributeValues); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
