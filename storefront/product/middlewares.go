package product

import (
	"net/http"
)

func (h *Handler) notFoundMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			h.renderNotFound(w, r, "Page not found")
			return
		}

		next.ServeHTTP(w, r)
	})
}
