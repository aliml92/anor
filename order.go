package anor

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
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
	ID              int64
	CartID          int64
	UserID          int64
	Status          OrderStatus
	TotalAmount     decimal.Decimal
	PaymentIntentID string
	ShippingAddress map[string]string
	BillingAddress  map[string]string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	OrderItems []*OrderItem
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

type OrderService interface {
	ConvertCartToOrder(ctx context.Context, cart Cart, piID string) (Order, error)
}

//type CreateOrderItemsParams struct {
//	OrderID           int64
//	VariantID         int64
//	Qty               int32
//	Price             decimal.Decimal
//	Thumbnail         string
//	ProductName       string
//	VariantAttributes []byte
//}
