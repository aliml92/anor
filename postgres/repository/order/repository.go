package order

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type CreateOrderItemsParams struct {
	OrderID           int64
	VariantID         int64
	Qty               int32
	Price             decimal.Decimal
	Thumbnail         string
	ProductName       string
	VariantAttributes []byte
}

type Repository interface {
	Querier
}

type repository struct {
	*Queries
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{Queries: New(db)}
}
