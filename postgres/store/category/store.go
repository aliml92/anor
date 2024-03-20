package category

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type categoryStore struct {
	*Queries
}

func NewStore(db *pgxpool.Pool) Store {
	return &categoryStore{Queries: New(db)}
}
