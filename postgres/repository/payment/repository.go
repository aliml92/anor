package payment

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Querier
}

type repository struct {
	*Queries
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{Queries: New(db)}
}
