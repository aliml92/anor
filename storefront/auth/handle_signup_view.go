package auth

import (
	"github.com/aliml92/anor/html/dtos/pages/signin"
	"github.com/aliml92/anor/html/dtos/partials"
	"net/http"
)

func (h *Handler) SignupView(w http.ResponseWriter, r *http.Request) {
	if isHXRequest(r) {
		h.view.Render(w, "pages/signup/content.gohtml", nil)
		return
	}

	ctx := r.Context()
	headerContent := partials.Header{}
	cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
	if cartId != 0 {
		guestCartItemCount, err := h.cartSvc.CountCartItems(ctx, cartId)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

		headerContent.CartNavItem = partials.CartNavItem{CartItemsCount: int(guestCartItemCount)}
	}

	sb := signin.Base{
		Header: headerContent,
	}

	h.view.Render(w, "pages/signup/base.gohtml", sb)
}
