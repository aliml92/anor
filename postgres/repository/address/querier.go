// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package address

import (
	"context"
)

type Querier interface {
	CreateAddress(ctx context.Context, userID int64, defaultFor NullAddressDefaultType, name string, addressLine1 string, addressLine2 *string, city string, stateProvince *string, postalCode *string, country *string) (*Address, error)
	GetAddressByID(ctx context.Context, id int64) (*Address, error)
	GetUserDefaultAddress(ctx context.Context, userID int64, defaultFor NullAddressDefaultType) (*Address, error)
	ListAddressesByUserID(ctx context.Context, userID int64, offset int32, limit int32) ([]*Address, error)
}

var _ Querier = (*Queries)(nil)
