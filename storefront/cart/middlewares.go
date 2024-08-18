package cart

import (
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/session"
	"net/http"
	"strconv"
)

func (h *Handler) VerifyCartItemOwnership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			h.clientError(w, err, http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		u := session.UserFromContext(ctx)
		if u.CartID == 0 {
			h.clientError(w, errors.New("no user found"), http.StatusBadRequest)
			return
		}

		isMyItem, err := h.cartService.IsMyItem(ctx, anor.IsMyCartItemParams{
			ItemID: itemID,
			CartID: u.CartID,
		})
		if err != nil {
			h.serverInternalError(w, err)
			return
		}
		if !isMyItem {
			h.clientError(w, errors.New("you don't have right access to this resource"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
