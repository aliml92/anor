package product

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/aliml92/anor"
	sb "github.com/aliml92/anor/postgres/sqlbuilder"
)

type Store interface {
	Querier
	GetProductsByLeafCategoryID(ctx context.Context, categoryID int32, p anor.GetProductsByCategoryParams) ([]ProductsByCategoryRow, error)
	GetProductsByLeafCategoryIDs(ctx context.Context, categoryID []int32, p anor.GetProductsByCategoryParams) ([]ProductsByCategoryRow, error)
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type productStore struct {
	*Queries
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &productStore{
		Queries: New(pool),
		pool:    pool,
	}
}

func (s productStore) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	if err := pgx.BeginFunc(ctx, s.pool, fn); err != nil {
		return err
	}
	return nil
}

type ProductsByCategoryRow struct {
	ID               int64
	Name             string
	Slug             string
	ThumbImg         string
	BasePrice        decimal.Decimal
	DiscountedAmount decimal.Decimal
	CurrencyCode     string
	Count            int64
}

func (s productStore) GetProductsByLeafCategoryID(
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
		"p.slug",
		"p.image_urls->>'0' as thumb_img",
		"pp.base_price",
		"pp.discounted_amount",
		"pp.currency_code",
		"count(p.id) over() as total_products",
	).
		From("products p").
		LeftJoin("product_pricing pp", "pp.product_id = p.id").
		Where().
		Eq("p.category_id", categoryID)

	brands := p.FilterOption.Brands
	if len(brands) > 0 {
		q.And()
		if len(brands) == 1 {
			q.Eq("p.brand", brands[0])
		} else {
			q.EqAny("p.brand", brands)
		}
	}

	priceRange := p.FilterOption.PriceRange
	// check if upper price is set
	if priceRange[1] != 0 {
		q.And()
		low := priceRange[0]
		high := priceRange[1]
		if low > 0 {
			q.Gt("pp.base_price", low)
			q.And()
			q.Lt("pp.base_price", high)
		} else {
			q.Lt("pp.base_price", high)
		}
	}

	switch p.SortOption {
	case anor.SortOptionBestMatch:
		// TODO: research on this case
	case anor.SortOptionPriceHighToLow:
		q.OrderBy(sb.M{"pp.price": "DESC"})
	case anor.SortOptionPriceLowToHigh:
		q.OrderBy(sb.M{"pp.price": "ASC"})
	case anor.SortOptionHighestRated:
		// TODO: implement later
	case anor.SortOptionNewArrivals:
		q.OrderBy(sb.M{"p.created_at": "DESC"})
	case anor.SortOptionBestSellers:
		// TODO: implement later
	}

	q.Offset(p.Offset)
	q.Limit(p.Limit)

	query, args := q.ToSql()
	println(query)

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
			&i.Slug,
			&i.ThumbImg,
			&i.BasePrice,
			&i.DiscountedAmount,
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

func (s productStore) GetProductsByLeafCategoryIDs(
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
		"p.slug",
		"p.image_urls->>'0' as thumb_img",
		"pp.base_price",
		"pp.discounted_amount",
		"pp.currency_code",
		"count(p.id) over() as total_products",
	).
		From("products p").
		LeftJoin("product_pricing pp", "pp.product_id = p.id").
		Where().
		EqAny("p.category_id", categoryIDs)

	brands := p.FilterOption.Brands
	if len(brands) > 0 {
		q.And()
		if len(brands) == 1 {
			q.Eq("p.brand", brands[0])
		} else {
			q.EqAny("p.brand", brands)
		}
	}

	priceRange := p.FilterOption.PriceRange
	// check if upper price is set
	if priceRange[1] != 0 {
		q.And()
		low := priceRange[0]
		high := priceRange[1]
		if low > 0 {
			q.Gt("pp.base_price", low)
			q.And()
			q.Lt("pp.base_price", high)
		} else {
			q.Lt("pp.base_price", high)
		}
	}

	switch p.SortOption {
	case anor.SortOptionBestMatch:
		// TODO: research on this case
	case anor.SortOptionPriceHighToLow:
		q.OrderBy(sb.M{"pp.price": "DESC"})
	case anor.SortOptionPriceLowToHigh:
		q.OrderBy(sb.M{"pp.price": "ASC"})
	case anor.SortOptionHighestRated:
		// TODO: implement later
	case anor.SortOptionNewArrivals:
		q.OrderBy(sb.M{"p.created_at": "DESC"})
	case anor.SortOptionBestSellers:
		// TODO: implement later
	}

	q.Offset(p.Offset)
	q.Limit(p.Limit)

	query, args := q.ToSql()
	println(query)

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
			&i.Slug,
			&i.ThumbImg,
			&i.BasePrice,
			&i.DiscountedAmount,
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
