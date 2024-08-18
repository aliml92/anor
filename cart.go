package anor

import (
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
	"time"
)

const DefaultCurrency = "USD"

type CartStatus string

const (
	CartStatusOpen      CartStatus = "Open"
	CartStatusCompleted CartStatus = "Completed"
	CartStatusMerged    CartStatus = "Merged"
	CartStatusExpired   CartStatus = "Expired"
	CartStatusAbandoned CartStatus = "Abandoned"
)

type Cart struct {
	ID           int64
	UserID       int64
	Status       CartStatus
	LastActivity time.Time
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	CartItems []*CartItem

	TotalAmount  decimal.Decimal
	CurrencyCode string
}

type CartItem struct {
	ID                int64
	CartID            int64
	VariantID         int64
	Qty               int32
	Price             decimal.Decimal
	Currency          string
	Thumbnail         string
	ProductName       string
	ProductPath       string
	VariantAttributes map[string]string
	IsRemoved         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time

	AvailableQty int32
}

type CartService interface {
	Create(ctx context.Context, params CartCreateParams) (Cart, error)
	Get(ctx context.Context, id int64, includeItems bool) (Cart, error)
	GetByUserID(ctx context.Context, userID int64, includeItems bool) (Cart, error)
	Update(ctx context.Context, id int64, params CartUpdateParams) error
	Merge(ctx context.Context, params CartMergeParams) (Cart, error)

	ListItems(ctx context.Context, params CartItemListParams) ([]CartItem, error)
	AddItem(ctx context.Context, params CartItemAddParams) (CartItem, error)
	UpdateItem(ctx context.Context, itemID int64, params CartItemUpdateParams) error
	DeleteItem(ctx context.Context, itemID int64) error

	CountItems(ctx context.Context, cartID int64) (int64, error)
	IsMyItem(ctx context.Context, params IsMyCartItemParams) (bool, error)
}

type CartCreateParams struct {
	UserID    int64
	ExpiresAt time.Time
}

type CartUpdateParams struct {
	Status CartStatus
}

type CartMergeParams struct {
	GuestCartID int64
	UserID      int64
}

type CartItemListParams struct {
	CartID   int64
	Page     int64
	PageSize int64
}

type CartItemAddParams struct {
	CartID    int64
	VariantID int64
	Qty       int
}

type CartItemUpdateParams struct {
	Qty int
}

type IsMyCartItemParams struct {
	ItemID int64
	CartID int64
}
