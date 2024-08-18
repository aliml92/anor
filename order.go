package anor

import (
	"context"
	"github.com/aliml92/anor/relation"
	"github.com/shopspring/decimal"
	"time"
)

type PaymentMethod string

const (
	PaymentMethodStripeCard      PaymentMethod = "StripeCard"
	PaymentMethodPaypal          PaymentMethod = "Paypal"
	PaymentMethodAnorInstallment PaymentMethod = "AnorInstallment"
	PaymentMethodPayOnDelivery   PaymentMethod = "PayOnDelivery"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "Pending"
	PaymentStatusPaid    PaymentStatus = "Paid"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "Pending"
	OrderStatusProcessing OrderStatus = "Processing"
	OrderStatusShipped    OrderStatus = "Shipped"
	OrderStatusDelivered  OrderStatus = "Delivered"
	OrderStatusCanceled   OrderStatus = "Canceled"
)

type Order struct {
	ID                int64
	UserID            int64
	CartID            int64
	PaymentMethod     PaymentMethod
	PaymentStatus     PaymentStatus
	Status            OrderStatus
	ShippingAddressID int64
	IsPickup          bool
	Amount            decimal.Decimal
	Currency          string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Related structs that can be eagerly loaded
	OrderItems        []*OrderItem
	ShippingAddress   Address
	StripeCardPayment StripeCardPayment
	DeliveryDate      time.Time
}

type OrderItem struct {
	ID                int64
	OrderID           int64
	VariantID         int64
	Qty               int32
	Price             decimal.Decimal
	Thumbnail         string
	ProductName       string
	VariantAttributes map[string]string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type OrderCreateParams struct {
	Cart              Cart
	PaymentMethod     PaymentMethod
	PaymentStatus     PaymentStatus
	ShippingAddressID int64
	IsPickup          bool
	Amount            decimal.Decimal
	Currency          string
}

type OrderListParams struct {
	UserID        int64
	WithRelations relation.Set
	Page          int
	PageSize      int
}

type OrderService interface {
	Create(ctx context.Context, params OrderCreateParams) (int64, error)
	Get(ctx context.Context, orderID int64, withItems bool) (Order, error)
	List(ctx context.Context, params OrderListParams) ([]Order, error)
	ListActive(ctx context.Context, params OrderListParams) ([]Order, error)
	ListUnpaid(ctx context.Context, params OrderListParams) ([]Order, error)
}
