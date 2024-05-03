package auth

import "net/http"

func (h *Handler) ForgotPasswordView(w http.ResponseWriter, r *http.Request) {
	h.renderView(w, r, http.StatusOK, "forgot-password.gohtml", nil)
}
