package checkout

import (
	"errors"
	"net/http"

	"github.com/aliml92/anor"
)

func (h *Handler) authInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// htmx requests do not require auth
		if isHXRequest(r) {
			next.ServeHTTP(w, r)
			return
		}
		var user *anor.User
		id := h.session.Auth.GetInt64(r.Context(), "authenticatedUserID")
		if id != 0 {
			// retrieve user by id from database
			ru, err := h.userSvc.GetUser(r.Context(), id) // retrieved user

			// check if error exists and it is not anor.ErrNotFound error
			if err != nil && !errors.Is(err, anor.ErrNotFound) {
				h.serverInternalError(w, err)
				return
			}

			// check if user is not found
			if errors.Is(err, anor.ErrNotFound) {
				user = nil
			} else {
				user = &ru
			}
		}

		ctx := anor.NewContextWithUser(r.Context(), user)
		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
