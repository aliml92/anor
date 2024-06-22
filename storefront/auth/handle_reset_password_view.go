package auth

import (
	"errors"
	"net/http"
)

func (h *Handler) ResetPasswordView(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		err := errors.New("invalid or expired reset password url. Request a new reset password link to proceed")
		h.clientError(w, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	ok, err := h.svc.VerifyResetPasswordToken(ctx, token)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}
	if !ok {
		err := errors.New("invalid or expired reset password url. Request a new reset password link to proceed")
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	if isHXRequest(r) {
		h.view.Render(w, "pages/reset-password/content.gohtml", nil)
		return
	}
	h.view.Render(w, "pages/reset-password/base.gohtml", nil)
}
