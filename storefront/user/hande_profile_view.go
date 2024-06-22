package user

import (
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/dtos/pages/profile"
	"github.com/aliml92/anor/html/dtos/partials"
	"net/http"
)

type ProfileViewData struct {
	User *anor.User
}

func (h *Handler) ProfileView(w http.ResponseWriter, r *http.Request) {
	if isHXRequest(r) {
		h.view.Render(w, "pages/profile/content.gohtml", nil)
		return
	}

	ctx := r.Context()
	u := anor.UserFromContext(ctx)
	headerContent := partials.Header{User: u}
	if u != nil {
		ac, err := h.userSvc.GetUserActivityCounts(ctx, u.ID)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

		headerContent.ActiveOrdersCount = ac.ActiveOrdersCount
		headerContent.WishlistItemsCount = ac.WishlistItemsCount
		headerContent.CartItemsCount = ac.CartItemsCount

	} else {
		cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
		if cartId != 0 {

			guestCartItemCount, err := h.cartSvc.CountCartItems(ctx, cartId)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}

			headerContent.CartItemsCount = int(guestCartItemCount)
		}
	}

	pb := profile.Base{
		Header: headerContent,
	}

	h.view.Render(w, "pages/profile/base.gohtml", pb)
}
