package auth

import "net/http"

func (h *Handler) SigninView(w http.ResponseWriter, r *http.Request) {
	if isHXRequest(r) {
		h.render.HTMX(w, http.StatusOK, "signin.gohtml", nil)
		return
	}

	h.render.HTML(w, http.StatusOK, "signin.gohtml", nil)
}
