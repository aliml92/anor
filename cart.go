package anor

import (
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
	"time"
)

type CartStatus string

const (
	CartStatusActive    CartStatus = "Active"
	CartStatusCheckout  CartStatus = "Checkout"
	CartStatusCompleted CartStatus = "Completed"
)

type Cart struct {
	ID             int64
	UserID         int64
	Status         CartStatus
	PIClientSecret string
	CreatedAt      time.Time
	DeletedAt      time.Time

	CartItems []*CartItem
	// extra data
	TotalAmount  decimal.Decimal
	CurrencyCode string
}

type CartItem struct {
	ID                int64             `json:"id"`
	CartID            int64             `json:"cart_id"`
	VariantID         int64             `json:"variant_id"`
	Qty               int32             `json:"qty"`
	Price             decimal.Decimal   `json:"price"`
	CurrencyCode      string            `json:"currency_code"`
	Thumbnail         string            `json:"thumbnail"`
	ProductName       string            `json:"product_name"`
	ProductPath       string            `json:"product_path"`
	VariantAttributes map[string]string `json:"variant_attributes"`
	IsRemoved         bool              `json:"is_removed"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`

	// extra data
	AvailableQty int32 `json:"available_qty"`
}

type AddCartItemParam struct {
	VariantID int64
	Qty       int
}

type UpdateCartItemParam struct {
	Qty int
}

type CartService interface {
	CreateCart(ctx context.Context, userID int64) (Cart, error)
	GetCart(ctx context.Context, userID int64, includeCartItems bool) (Cart, error)
	GetGuestCartItems(ctx context.Context, cartID int64) ([]*CartItem, error)
	AddCartItem(ctx context.Context, cartID int64, p AddCartItemParam) (CartItem, error)
	UpdateCartItem(ctx context.Context, cartItemID int64, p UpdateCartItemParam) error
	DeleteCartItem(ctx context.Context, cartItemID int64) error
	UpdateCart(ctx context.Context, c Cart) error

	CountCartItems(ctx context.Context, cartID int64) (int64, error)
	IsCartItemOwner(ctx context.Context, userID int64, cartItemId int64) (bool, error)
	IsGuestCartItemOwner(ctx context.Context, cartID int64, cartItemId int64) (bool, error)
}
