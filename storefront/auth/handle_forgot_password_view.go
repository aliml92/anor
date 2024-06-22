package auth

import "net/http"

func (h *Handler) ForgotPasswordView(w http.ResponseWriter, r *http.Request) {
	if isHXRequest(r) {
		h.view.Render(w, "pages/forgot-password/content.gohtml", nil)
		return
	}
	h.view.Render(w, "pages/forgot-password/base.gohtml", nil)
}
