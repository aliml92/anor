package cart

import (
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/dtos/pages/cart"
	"github.com/aliml92/anor/html/dtos/pages/cart/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/shopspring/decimal"
	"net/http"
)

func (h *Handler) GetCartView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := anor.UserFromContext(ctx)

	var (
		cartItems []*anor.CartItem
		err       error
	)
	if u != nil {
		c, err := h.cartSvc.GetCart(ctx, u.ID, true)
		if err != nil && !errors.Is(err, anor.ErrNotFound) {
			h.serverInternalError(w, err)
			return
		}
		cartItems = c.CartItems

	} else {
		cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
		if cartId != 0 {
			cartItems, err = h.cartSvc.GetGuestCartItems(ctx, cartId)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}
		}
	}

	cc := cart.Content{
		CartItems: components.CartItems{CartItems: cartItems},
		CartSummary: components.CartSummary{
			CartItemsCount: len(cartItems),
			TotalAmount:    getTotalPrice(cartItems),
			CurrencyCode:   getCurrency(cartItems),
		},
	}

	if isHXRequest(r) {
		h.view.Render(w, "pages/cart/content.gohtml", cc)
		return
	}

	hxSwapOOB := r.Header.Get("Hx-Swap-OOB")
	cb := &cart.Base{
		Header: partials.Header{
			User:            u,
			OrdersNavItem:   partials.OrdersNavItem{},   // TODO: fill in orders data
			WishlistNavItem: partials.WishlistNavItem{}, // TODO: fill in wishlist data
			CartNavItem: partials.CartNavItem{
				HxSwapOOB:      hxSwapOOB,
				CartItemsCount: len(cartItems),
			},
		},
		Content: cc,
	}
	h.view.Render(w, "pages/cart/base.gohtml", cb)
}

// getTotalPrice calculates total price of cart, it is assumed that
// all cart items have same currency code
func getTotalPrice(cartItems []*anor.CartItem) decimal.Decimal {
	var totalPrice decimal.Decimal
	for _, cartItem := range cartItems {
		cost := decimal.NewFromInt32(cartItem.Qty).Mul(cartItem.Price)
		totalPrice = totalPrice.Add(cost)
	}
	return totalPrice
}

// getCurrency makes a naive assumption about cart currency code, it is assumed that
// all cart items have same currency code
func getCurrency(cartItems []*anor.CartItem) string {
	if len(cartItems) > 0 {
		return cartItems[0].CurrencyCode
	}

	return "USD"
}
