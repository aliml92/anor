package components

import (
	"github.com/aliml92/anor"
	"github.com/shopspring/decimal"
)

type CartItems struct {
	CartItems []*anor.CartItem
}

func (CartItems) GetTemplateFilename() string {
	return "cart_items.gohtml"
}

type CartSummary struct {
	HxSwapOOB      bool
	CartItemsCount int
	TotalAmount    decimal.Decimal
	Currency       string
}

func (CartSummary) GetTemplateFilename() string {
	return "cart_summary.gohtml"
}
