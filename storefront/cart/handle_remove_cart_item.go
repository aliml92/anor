package cart

import (
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
	cartCompontents "github.com/aliml92/anor/html/templates/pages/cart/components"
	headerComponents "github.com/aliml92/anor/html/templates/shared/header/components"
	"github.com/aliml92/anor/session"
	"net/http"
	"strconv"
)

func (h *Handler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.cartService.DeleteItem(ctx, itemID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	u := session.UserFromContext(ctx)
	c, err := h.cartService.Get(ctx, u.CartID, true)
	if err != nil && !errors.Is(err, anor.ErrCartNotFound) {
		h.serverInternalError(w, err)
		return
	}

	w.Header().Add("HX-Trigger-After-Settle", `{"anor:showToast":"item removed successfully"}`)

	itemCount := len(c.CartItems)

	cnv := headerComponents.CartNavItem{
		HxSwapOOB:      true,
		CartItemsCount: itemCount,
	}

	comps := []html.Component{
		{"shared/header/components/cart_nav_item.gohtml": cnv},
	}

	if itemCount == 0 {
		w.Header().Add("HX-Retarget", "#cart-content")
		comps = append(comps, html.Component{
			"pages/cart/components/cart_empty.gohtml": nil,
		})
	} else {
		cs := cartCompontents.CartSummary{
			HxSwapOOB:      true,
			CartItemsCount: itemCount,
			TotalAmount:    calculateTotalAmount(c.CartItems),
			Currency:       anor.DefaultCurrency,
		}
		comps = append(comps, html.Component{
			"pages/cart/components/cart_summary.gohtml": cs,
		})
	}

	h.view.RenderComponents(w, comps)
}
