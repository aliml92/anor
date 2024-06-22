package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Querier
	ExecTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	WithTx(tx pgx.Tx) *Queries
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

func (s repository) WithTx(tx pgx.Tx) *Queries {
	return s.Queries.WithTx(tx)
}

func (s repository) ExecTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	if err := pgx.BeginFunc(ctx, s.pool, fn); err != nil {
		return err
	}

	return nil
}
