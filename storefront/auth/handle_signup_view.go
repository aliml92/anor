package auth

import (
	"github.com/aliml92/anor/html/templates/pages/auth/signup"
	"net/http"
)

func (h *Handler) SignupView(w http.ResponseWriter, r *http.Request) {
	c := signup.Content{}
	h.Render(w, r, "pages/auth/signup/content.gohtml", c)
}
