package middlewares

import (
	"github.com/aliml92/anor/session"
	"net/http"
	"net/url"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u := session.UserFromContext(ctx)
		if !u.IsAuth {
			redirectURL, err := url.Parse("/auth/signin")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			query := redirectURL.Query()
			query.Set("redirect_url", r.URL.String())
			redirectURL.RawQuery = query.Encode()
			http.Redirect(w, r, redirectURL.String(), http.StatusSeeOther)
			return
		}

		//user, err := h.userService.getUser(r.Context(), u.UserID)
		//if err != nil {
		//	if errors.Is(err, anor.ErrNotFound) {
		//		if err := h.session.Destroy(ctx); err != nil {
		//			h.logger.Error("Failed to destroy session", "error", err)
		//			h.serverInternalError(w, err)
		//			return
		//		}
		//		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		//	} else {
		//		h.serverInternalError(w, err)
		//	}
		//	return
		//}

		//ctx = anor.NewContextWithUser(ctx, &user)
		//next.ServeHTTP(w, r.WithContext(ctx))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
