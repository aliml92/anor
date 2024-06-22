package product

import (
	"errors"
	not_found "github.com/aliml92/anor/html/dtos/pages/not-found"
	"github.com/aliml92/anor/html/dtos/partials"
	"net/http"

	"github.com/aliml92/anor"
)

func (h *Handler) authInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *anor.User
		id := h.session.Auth.GetInt64(r.Context(), "authenticatedUserID")
		if id != 0 {
			// retrieve user by id from database
			ru, err := h.userSvc.GetUser(r.Context(), id) // retrieved user
			if errors.Is(err, anor.ErrNotFound) {
				user = nil
			} else if err != nil {
				h.serverInternalError(w, err)
				return
			} else {
				user = &ru
			}
		}

		ctx := anor.NewContextWithUser(r.Context(), user)
		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) notFoundMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			nfc := not_found.Content{Message: "Page not found"}
			if isHXRequest(r) {
				h.view.Render(w, "pages/not-found/content.gohtml", nfc)
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

			nfb := not_found.Base{
				Header:  headerContent,
				Content: nfc,
			}

			h.view.Render(w, "pages/not-found/base.gohtml", nfb)
			return
		}

		next.ServeHTTP(w, r)
	})
}
