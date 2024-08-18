package anor

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
)

type StripeCardPayment struct {
	ID               int64
	OrderID          int64
	UserID           int64
	BillingAddressID int64
	PaymentIntentID  string
	PaymentMethodID  string
	Amount           decimal.Decimal
	Currency         string
	Status           string
	ClientSecret     string
	LastError        string
	CardLast4        string
	CardBrand        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type StripePaymentCreateParams struct {
	OrderID          int64
	UserID           int64
	BillingAddressID int64
	PaymentIntentID  string
	PaymentMethodID  string
	Amount           decimal.Decimal
	Currency         string
	Status           string
	ClientSecret     string
	LastError        string
	CardLast4        string
	CardBrand        string
	CreatedAt        time.Time
}

type StripePaymentService interface {
	GetByOrderID(ctx context.Context, orderID int64) (StripeCardPayment, error)
	Create(ctx context.Context, params StripePaymentCreateParams) error
}
