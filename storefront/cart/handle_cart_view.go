package cart

import (
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/pages/cart"
	"github.com/aliml92/anor/html/templates/pages/cart/components"
	"github.com/aliml92/anor/session"
	"github.com/shopspring/decimal"
	"net/http"
)

func (h *Handler) CartView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := session.UserFromContext(ctx)
	if u.CartID == 0 {
		h.Render(w, r, "pages/cart/content.gohtml", cart.Content{})
		return
	}

	c, err := h.cartService.Get(ctx, u.CartID, true)
	if err != nil && !errors.Is(err, anor.ErrCartNotFound) {
		h.serverInternalError(w, err)
		return
	}

	cc := cart.Content{
		CartItems: components.CartItems{CartItems: c.CartItems},
		CartSummary: components.CartSummary{
			CartItemsCount: len(c.CartItems),
			TotalAmount:    calculateTotalAmount(c.CartItems),
			Currency:       anor.DefaultCurrency,
		},
	}

	h.Render(w, r, "pages/cart/content.gohtml", cc)
}

func calculateTotalAmount(cartItems []*anor.CartItem) decimal.Decimal {
	var totalPrice decimal.Decimal
	for _, cartItem := range cartItems {
		cost := decimal.NewFromInt32(cartItem.Qty).Mul(cartItem.Price)
		totalPrice = totalPrice.Add(cost)
	}
	return totalPrice
}
