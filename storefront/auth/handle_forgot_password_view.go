package auth

import (
	forgotpassword "github.com/aliml92/anor/html/templates/pages/auth/forgot_password"
	"net/http"
)

func (h *Handler) ForgotPasswordView(w http.ResponseWriter, r *http.Request) {
	fc := forgotpassword.Content{}
	h.Render(w, r, "pages/auth/forgot_password/content.gohtml", fc)
}
