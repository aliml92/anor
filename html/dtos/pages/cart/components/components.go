package components

import (
	"github.com/aliml92/anor"
	"github.com/shopspring/decimal"
)

type CartItems struct {
	CartItems []*anor.CartItem
}

type CartSummary struct {
	HxSwapOOB      string
	CartItemsCount int
	TotalAmount    decimal.Decimal
	CurrencyCode   string
}
