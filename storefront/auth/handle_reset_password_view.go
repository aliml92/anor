package auth

import (
	"errors"
	resetpassword "github.com/aliml92/anor/html/templates/pages/auth/reset_password"
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
	ok, err := h.authService.VerifyResetPasswordToken(ctx, token)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}
	if !ok {
		err := errors.New("invalid or expired reset password url. Request a new reset password link to proceed")
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	rc := resetpassword.Content{}
	h.Render(w, r, "pages/auth/reset_password/content.gohtml", rc)
}
