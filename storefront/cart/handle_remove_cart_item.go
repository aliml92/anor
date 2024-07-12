package cart

import (
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
	"github.com/aliml92/anor/html/dtos/pages/cart/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"net/http"
	"strconv"
)

func (h *Handler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	cartItemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.cartSvc.DeleteCartItem(ctx, cartItemID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	u := anor.UserFromContext(ctx)
	c, err := h.cartSvc.GetCart(ctx, u.ID, true)
	if err != nil && !errors.Is(err, anor.ErrNotFound) {
		h.serverInternalError(w, err)
		return
	}

	w.Header().Add("HX-Trigger-After-Settle", `{"anor:showToast":"item removed successfully"}`)

	hxSwapOOB := r.Header.Get("Hx-Swap-OOB")
	itemCount := len(c.CartItems)

	cs := components.CartSummary{
		HxSwapOOB:      hxSwapOOB,
		CartItemsCount: itemCount,
		TotalAmount:    getTotalPrice(c.CartItems),
		CurrencyCode:   getCurrency(c.CartItems),
	}
	cnv := partials.CartNavItem{
		HxSwapOOB:      hxSwapOOB,
		CartItemsCount: itemCount,
	}

	comps := []html.Component{
		{"pages/cart/components/cart-summary.gohtml": cs},
		{"partials/header/cart-nav-item.gohtml": cnv},
	}

	if itemCount == 0 {
		comps = append(comps, html.Component{
			"pages/cart/components/cart-empty.gohtml": nil,
		})
	}
	h.view.RenderComponents(w, comps)
}
