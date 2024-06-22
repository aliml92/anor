package auth

import (
	"github.com/aliml92/anor/html/dtos/pages/signin"
	"github.com/aliml92/anor/html/dtos/partials"
	"net/http"
)

func (h *Handler) SigninView(w http.ResponseWriter, r *http.Request) {
	sc := signin.Content{}
	if isHXRequest(r) {
		h.view.Render(w, "pages/signin/content.gohtml", sc)
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

		headerContent.CartItemsCount = int(guestCartItemCount)
	}

	sb := signin.Base{
		Header:  headerContent,
		Content: sc,
	}

	h.view.Render(w, "pages/signin/base.gohtml", sb)
}
