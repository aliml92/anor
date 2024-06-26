// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package cart

import (
	"context"

	"github.com/shopspring/decimal"
)

type Querier interface {
	AddCartItem(ctx context.Context, cartID int64, variantID int64, qty int32, price decimal.Decimal, currencyCode string, thumbnail string, productName string, productPath string, variantAttributes []byte) (*CartItem, error)
	CartItemExistsByCartID(ctx context.Context, iD int64, cartID int64) (bool, error)
	CartItemExistsByUserID(ctx context.Context, iD int64, iD_2 int64) (bool, error)
	CountCartItems(ctx context.Context, userID *int64) (int64, error)
	CountCartItemsByCartID(ctx context.Context, cartID int64) (int64, error)
	CreateCart(ctx context.Context, userID *int64) (*Cart, error)
	CreateGuestCart(ctx context.Context) (*Cart, error)
	DeleteCartItem(ctx context.Context, id int64) error
	GetCartByID(ctx context.Context, id int64) (*Cart, error)
	GetCartByUserID(ctx context.Context, userID *int64) (*Cart, error)
	GetCartItemsByCartID(ctx context.Context, cartID int64) ([]*CartItem, error)
	GetCartItemsByUserID(ctx context.Context, userID *int64) ([]*CartItem, error)
	IncrementCartItemQty(ctx context.Context, cartID int64, variantID int64, qty int32) (*CartItem, error)
	UpdateCartClientSecret(ctx context.Context, iD int64, userID *int64, piClientSecret *string) error
	UpdateCartItemQty(ctx context.Context, iD int64, qty int32) error
}

var _ Querier = (*Queries)(nil)
