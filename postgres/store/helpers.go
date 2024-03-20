package store

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func Nullable[T any](row *T, err error) (*T, error) {
	if err == nil {
		return row, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
