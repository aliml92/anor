package user

import (
	profilepage "github.com/aliml92/anor/html/templates/pages/profile"
	"net/http"
)

func (h *Handler) ProfileView(w http.ResponseWriter, r *http.Request) {
	var c profilepage.Content
	h.Render(w, r, "pages/profile/content.gohtml", c)
}
