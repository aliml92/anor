package auth

import "net/http"

func (h *Handler) SignupView(w http.ResponseWriter, r *http.Request) {
	h.renderView(w, r, http.StatusOK, "signup.gohtml", nil)
}
