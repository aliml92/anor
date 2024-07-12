package cart

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aliml92/anor"
)

func (h *Handler) isAuthorizedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cartItemIDStr := r.PathValue("id")
		cartItemID, err := strconv.ParseInt(cartItemIDStr, 10, 64)
		if err != nil {
			h.clientError(w, err, http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		u := anor.UserFromContext(ctx)
		if u != nil {
			isOwner, err := h.cartSvc.IsCartItemOwner(ctx, u.ID, cartItemID)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}
			if !isOwner {
				h.clientError(w, errors.New("you don't have right access to this resource"), http.StatusForbidden)
				return
			}
		} else {
			isOwner, err := h.cartSvc.IsGuestCartItemOwner(ctx, u.ID, cartItemID)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}
			if !isOwner {
				h.clientError(w, errors.New("you don't have right access to this resource"), http.StatusForbidden)
				return
			}
		}

		// Call the next handler with the new context
		next.ServeHTTP(w, r)
	})
}
