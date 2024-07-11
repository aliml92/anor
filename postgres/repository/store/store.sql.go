// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: store.sql

package store

import (
	"context"
)

const createStore = `-- name: CreateStore :one
INSERT INTO stores (
    handle, user_id, name, description
) VALUES (
    $1, $2, $3, $4
) RETURNING id
`

func (q *Queries) CreateStore(ctx context.Context, handle string, userID int64, name string, description string) (int32, error) {
	row := q.db.QueryRow(ctx, createStore,
		handle,
		userID,
		name,
		description,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}