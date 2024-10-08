// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: product.sql

package product

import (
	"context"
)

const collectCartData = `-- name: CollectCartData :one
SELECT
    p.id AS product_id,
    p.name AS product_name,
    p.handle AS product_handle,
    p.image_urls AS product_image_urls,
    pv.sku,
    pv.qty,
    pv.is_custom_priced,
    pv.image_identifiers,
    COALESCE(
                    json_object_agg(
                    pa.attribute,
                    pva.attribute_value
                                   ) FILTER (WHERE pa.id IS NOT NULL),
                    '{}'
    )::json AS variant_attributes
FROM
    product_variants pv
        JOIN
    products p ON pv.product_id = p.id
        LEFT JOIN
    product_variant_attributes pva ON pv.id = pva.variant_id
        LEFT JOIN
    product_attributes pa ON pva.product_attribute_id = pa.id
WHERE
    pv.id = $1
GROUP BY
    p.id, p.name, p.image_urls, pv.sku, pv.qty, pv.is_custom_priced, pv.image_identifiers
`

type CollectCartDataRow struct {
	ProductID         int64
	ProductName       string
	ProductHandle     string
	ProductImageUrls  ImageUrls
	Sku               string
	Qty               int32
	IsCustomPriced    bool
	ImageIdentifiers  []int16
	VariantAttributes []byte
}

func (q *Queries) CollectCartData(ctx context.Context, id int64) (*CollectCartDataRow, error) {
	row := q.db.QueryRow(ctx, collectCartData, id)
	var i CollectCartDataRow
	err := row.Scan(
		&i.ProductID,
		&i.ProductName,
		&i.ProductHandle,
		&i.ProductImageUrls,
		&i.Sku,
		&i.Qty,
		&i.IsCustomPriced,
		&i.ImageIdentifiers,
		&i.VariantAttributes,
	)
	return &i, err
}

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
    store_id, category_id, name, brand, handle, short_information, image_urls, specifications, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, store_id, category_id, name, brand, handle, image_urls, short_information, specifications, status, created_at, updated_at
`

func (q *Queries) CreateProduct(ctx context.Context, storeID int32, categoryID int32, name string, brand *string, handle string, shortInformation []string, imageUrls ImageUrls, specifications Specifications, status ProductStatus) (*Product, error) {
	row := q.db.QueryRow(ctx, createProduct,
		storeID,
		categoryID,
		name,
		brand,
		handle,
		shortInformation,
		imageUrls,
		specifications,
		status,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.StoreID,
		&i.CategoryID,
		&i.Name,
		&i.Brand,
		&i.Handle,
		&i.ImageUrls,
		&i.ShortInformation,
		&i.Specifications,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getProductBrandsByCategoryID = `-- name: GetProductBrandsByCategoryID :many
SELECT array_agg(DISTINCT p.brand ORDER BY p.brand)::text[] AS brands FROM products p
WHERE p.category_id = $1
`

func (q *Queries) GetProductBrandsByCategoryID(ctx context.Context, categoryID int32) ([][]string, error) {
	rows, err := q.db.Query(ctx, getProductBrandsByCategoryID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][]string
	for rows.Next() {
		var brands []string
		if err := rows.Scan(&brands); err != nil {
			return nil, err
		}
		items = append(items, brands)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductBrandsByCategoryIDs = `-- name: GetProductBrandsByCategoryIDs :many
SELECT array_agg(DISTINCT p.brand ORDER BY p.brand)::text[] AS brands FROM products p
WHERE p.category_id = ANY($1::INT[])
`

func (q *Queries) GetProductBrandsByCategoryIDs(ctx context.Context, leafcategoryids []int32) ([][]string, error) {
	rows, err := q.db.Query(ctx, getProductBrandsByCategoryIDs, leafcategoryids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][]string
	for rows.Next() {
		var brands []string
		if err := rows.Scan(&brands); err != nil {
			return nil, err
		}
		items = append(items, brands)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductByID = `-- name: GetProductByID :one
SELECT id, store_id, category_id, name, brand, handle, image_urls, short_information, specifications, status, created_at, updated_at FROM products
WHERE id = $1
`

func (q *Queries) GetProductByID(ctx context.Context, id int64) (*Product, error) {
	row := q.db.QueryRow(ctx, getProductByID, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.StoreID,
		&i.CategoryID,
		&i.Name,
		&i.Brand,
		&i.Handle,
		&i.ImageUrls,
		&i.ShortInformation,
		&i.Specifications,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getProductImageURLsByID = `-- name: GetProductImageURLsByID :one
SELECT image_urls
FROM products
WHERE id = $1
`

func (q *Queries) GetProductImageURLsByID(ctx context.Context, id int64) (ImageUrls, error) {
	row := q.db.QueryRow(ctx, getProductImageURLsByID, id)
	var image_urls ImageUrls
	err := row.Scan(&image_urls)
	return image_urls, err
}

const getProductWithStoreAndPricingByID = `-- name: GetProductWithStoreAndPricingByID :one
SELECT
    p.id, p.store_id, p.category_id, p.name, p.brand, p.handle, p.image_urls, p.short_information, p.specifications, p.status, p.created_at, p.updated_at,
    s.id, s.handle, s.user_id, s.name, s.description, s.created_at, s.updated_at,
    pp.product_id, pp.base_price, pp.currency_code, pp.discount, pp.discounted_price, pp.is_on_sale
FROM
    products p
        LEFT JOIN
            product_pricing pp ON pp.product_id = p.id
        LEFT JOIN
            stores s ON s.id = p.store_id
WHERE
    p.id = $1
`

type GetProductWithStoreAndPricingByIDRow struct {
	Product        Product
	Store          Store
	ProductPricing ProductPricing
}

func (q *Queries) GetProductWithStoreAndPricingByID(ctx context.Context, id int64) (*GetProductWithStoreAndPricingByIDRow, error) {
	row := q.db.QueryRow(ctx, getProductWithStoreAndPricingByID, id)
	var i GetProductWithStoreAndPricingByIDRow
	err := row.Scan(
		&i.Product.ID,
		&i.Product.StoreID,
		&i.Product.CategoryID,
		&i.Product.Name,
		&i.Product.Brand,
		&i.Product.Handle,
		&i.Product.ImageUrls,
		&i.Product.ShortInformation,
		&i.Product.Specifications,
		&i.Product.Status,
		&i.Product.CreatedAt,
		&i.Product.UpdatedAt,
		&i.Store.ID,
		&i.Store.Handle,
		&i.Store.UserID,
		&i.Store.Name,
		&i.Store.Description,
		&i.Store.CreatedAt,
		&i.Store.UpdatedAt,
		&i.ProductPricing.ProductID,
		&i.ProductPricing.BasePrice,
		&i.ProductPricing.CurrencyCode,
		&i.ProductPricing.Discount,
		&i.ProductPricing.DiscountedPrice,
		&i.ProductPricing.IsOnSale,
	)
	return &i, err
}

const listPopularProducts = `-- name: ListPopularProducts :many
SELECT p.id, p.store_id, p.category_id, p.name, p.brand, p.handle, p.image_urls, p.short_information, p.specifications, p.status, p.created_at, p.updated_at, pp.product_id, pp.base_price, pp.currency_code, pp.discount, pp.discounted_price, pp.is_on_sale FROM products p
LEFT JOIN product_pricing pp ON p.id = pp.product_id
WHERE p.id = ANY($3::BIGINT[])
ORDER BY p.created_at DESC
LIMIT $1
    OFFSET $2
`

type ListPopularProductsRow struct {
	Product        Product
	ProductPricing ProductPricing
}

func (q *Queries) ListPopularProducts(ctx context.Context, limit int32, offset int32, productIDs []int64) ([]*ListPopularProductsRow, error) {
	rows, err := q.db.Query(ctx, listPopularProducts, limit, offset, productIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListPopularProductsRow
	for rows.Next() {
		var i ListPopularProductsRow
		if err := rows.Scan(
			&i.Product.ID,
			&i.Product.StoreID,
			&i.Product.CategoryID,
			&i.Product.Name,
			&i.Product.Brand,
			&i.Product.Handle,
			&i.Product.ImageUrls,
			&i.Product.ShortInformation,
			&i.Product.Specifications,
			&i.Product.Status,
			&i.Product.CreatedAt,
			&i.Product.UpdatedAt,
			&i.ProductPricing.ProductID,
			&i.ProductPricing.BasePrice,
			&i.ProductPricing.CurrencyCode,
			&i.ProductPricing.Discount,
			&i.ProductPricing.DiscountedPrice,
			&i.ProductPricing.IsOnSale,
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
