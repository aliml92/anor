package sellerstore

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type sellerStore struct {
	*Queries
}

func NewSellerStore(db *pgxpool.Pool) Store {
	return &sellerStore{Queries: New(db)}
}
