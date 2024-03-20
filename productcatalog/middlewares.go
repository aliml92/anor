package productcatalog

import (
	"context"
	"errors"
	"net/http"

	"github.com/aliml92/anor"
)

func (h *Handler) authInjector(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *anor.User
		id := h.session.GetInt64(r.Context(), "authenticatedUserID")
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
				*user = ru
			}
		}

		ctx := context.WithValue(r.Context(), "auth", user)

		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (h *Handler) NotFound(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			if isHXRequest(r) {
				h.render.HTMX(w, http.StatusOK, "not-found.tmpl", nil)
				return
			}

			ctx := r.Context()
			user := ctx.Value("auth").(*anor.User)

			h.render.HTML(w, http.StatusOK, "not-found.tmpl", HomeViewData{User: user})
			return
		}

		next.ServeHTTP(w, r)
	}
}
