// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package store

import (
	"context"
)

type Querier interface {
	CreateStore(ctx context.Context, handle string, userID int64, name string, description string) (int32, error)
}

var _ Querier = (*Queries)(nil)
