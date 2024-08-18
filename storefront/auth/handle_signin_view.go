package auth

import (
	"github.com/aliml92/anor/html/templates/pages/auth/signin"
	"net/http"
)

func (h *Handler) SigninView(w http.ResponseWriter, r *http.Request) {
	sc := signin.Content{}
	h.Render(w, r, "pages/auth/signin/content.gohtml", sc)
}
