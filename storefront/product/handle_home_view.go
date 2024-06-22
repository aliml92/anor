package product

import (
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/dtos/pages/home"
	"github.com/aliml92/anor/html/dtos/pages/home/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"net/http"
)


func (h *Handler) HomeView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	newArrivals, err := h.productSvc.GetNewArrivals(ctx, 10)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	hc := home.Content{
		Featured: components.Featured{
			Products: nil, // TODO: get featured products for carousel
		},
		NewArrivals: components.Collection{
			Products: newArrivals,
		},
		Popular: components.Collection{
			Products: nil, // TODO: get popular products
		},
	}
	if isHXRequest(r) {
		h.view.Render(w, "pages/home/content.gohtml", hc)
		return
	}

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

	hb := home.Base{
		Header:  headerContent,
		Content: hc,
	}

	h.view.Render(w, "pages/home/base.gohtml", hb)
}
