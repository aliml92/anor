package product

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/aliml92/anor"
	sb "github.com/aliml92/anor/postgres/sqlbuilder"
)

type Repository interface {
	Querier
	GetProductsByLeafCategoryID(ctx context.Context, categoryID int32, p anor.GetProductsByCategoryParams) ([]ProductsByCategoryRow, error)
	GetProductsByLeafCategoryIDs(ctx context.Context, categoryID []int32, p anor.GetProductsByCategoryParams) ([]ProductsByCategoryRow, error)
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type repository struct {
	*Queries
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{
		Queries: New(pool),
		pool:    pool,
	}
}

func (s repository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	if err := pgx.BeginFunc(ctx, s.pool, fn); err != nil {
		return err
	}
	return nil
}

type ProductsByCategoryRow struct {
	ID              int64
	Name            string
	Handle          string
	ImageUrls       ImageUrls
	BasePrice       decimal.Decimal
	Discount        decimal.Decimal
	DiscountedPrice decimal.Decimal
	CurrencyCode    string
	Count           int64
}

func (s repository) GetProductsByLeafCategoryID(
	ctx context.Context,
	categoryID int32,
	p anor.GetProductsByCategoryParams,
) ([]ProductsByCategoryRow, error) {
	// TODO: filtering by color, size, rating

	// apply price range and brand if exists
	q := sb.NewSqlBuilder()
	q.Select(
		"p.id",
		"p.name",
		"p.handle",
		"p.image_urls",
		"pp.base_price",
		"pp.discount",
		"pp.discounted_price",
		"pp.currency_code",
		"count(p.id) over() as total_products",
	).
		From("products p").
		LeftJoin("product_pricing pp", "pp.product_id = p.id").
		Where().
		Eq("p.category_id", categoryID)

	brands := p.Filter.Brands
	if len(brands) > 0 {
		q.And()
		if len(brands) == 1 {
			q.Eq("p.brand", brands[0])
		} else {
			q.EqAny("p.brand", brands)
		}
	}

	priceFrom := p.Filter.PriceFrom
	if !priceFrom.IsZero() {
		q.And()
		q.Ge("pp.discounted_price", priceFrom)
	}

	priceTo := p.Filter.PriceTo
	if !priceTo.IsZero() {
		q.And()
		q.Le("pp.discounted_price", priceTo)
	}

	switch p.Sort {
	case anor.SortParamPopular:
		// TODO: research on this case
	case anor.SortParamPriceHighToLow:
		q.OrderBy(sb.M{"pp.discounted_price": "DESC"})
	case anor.SortParamPriceLowToHigh:
		q.OrderBy(sb.M{"pp.discounted_price": "ASC"})
	case anor.SortParamHighestRated:
		// TODO: implement later
	case anor.SortParamNewArrivals:
		q.OrderBy(sb.M{"p.created_at": "DESC"})
	case anor.SortParamBestSellers:
		// TODO: implement later
	}

	q.Offset((p.Paging.Page - 1) * p.Paging.PageSize)
	q.Limit(p.Paging.PageSize)

	query, args := q.ToSql()

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductsByCategoryRow
	for rows.Next() {
		var i ProductsByCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Handle,
			&i.ImageUrls,
			&i.BasePrice,
			&i.Discount,
			&i.DiscountedPrice,
			&i.CurrencyCode,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s repository) GetProductsByLeafCategoryIDs(
	ctx context.Context,
	categoryIDs []int32,
	p anor.GetProductsByCategoryParams,
) ([]ProductsByCategoryRow, error) {
	// TODO: filtering by color, size, rating

	// apply price range and brand if exists
	q := sb.NewSqlBuilder()
	q = q.Select(
		"p.id",
		"p.name",
		"p.handle",
		"p.image_urls",
		"pp.base_price",
		"pp.discount",
		"pp.discounted_price",
		"pp.currency_code",
		"count(p.id) over() as total_products",
	).
		From("products p").
		LeftJoin("product_pricing pp", "pp.product_id = p.id").
		Where().
		EqAny("p.category_id", categoryIDs)

	brands := p.Filter.Brands
	if len(brands) > 0 {
		q.And()
		if len(brands) == 1 {
			q.Eq("p.brand", brands[0])
		} else {
			q.EqAny("p.brand", brands)
		}
	}

	priceFrom := p.Filter.PriceFrom
	if !priceFrom.IsZero() {
		q.And()
		q.Ge("pp.discounted_price", priceFrom)
	}

	priceTo := p.Filter.PriceTo
	if !priceTo.IsZero() {
		q.And()
		q.Le("pp.discounted_price", priceTo)
	}

	switch p.Sort {
	case anor.SortParamPopular:
		// TODO: research on this case
	case anor.SortParamPriceHighToLow:
		q.OrderBy(sb.M{"pp.discounted_price": "DESC"})
	case anor.SortParamPriceLowToHigh:
		q.OrderBy(sb.M{"pp.discounted_price": "ASC"})
	case anor.SortParamHighestRated:
		// TODO: implement later
	case anor.SortParamNewArrivals:
		q.OrderBy(sb.M{"p.created_at": "DESC"})
	case anor.SortParamBestSellers:
		// TODO: implement later
	}

	q.Offset((p.Paging.Page - 1) * p.Paging.PageSize)
	q.Limit(p.Paging.PageSize)

	query, args := q.ToSql()

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductsByCategoryRow
	for rows.Next() {
		var i ProductsByCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Handle,
			&i.ImageUrls,
			&i.BasePrice,
			&i.Discount,
			&i.DiscountedPrice,
			&i.CurrencyCode,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
