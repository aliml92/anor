package auth

import (
	"github.com/aliml92/anor/session"
	"net/http"
)

func (h *Handler) RedirectAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u := session.UserFromContext(ctx)
		if u.IsAuth {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
