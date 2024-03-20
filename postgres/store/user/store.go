package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	WithTx(tx pgx.Tx) *Queries
}

type store struct {
	*Queries
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &store{
		Queries: New(pool),
		pool:    pool,
	}
}

func (s store) WithTx(tx pgx.Tx) *Queries {
	return s.Queries.WithTx(tx)
}

func (s store) ExecTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	if err := pgx.BeginFunc(ctx, s.pool, fn); err != nil {
		return err
	}

	return nil
}
