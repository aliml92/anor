package user

import (
	"github.com/aliml92/anor"
	"net/http"
)

type ProfileViewData struct {
	User *anor.User
}

func (h *Handler) ProfileView(w http.ResponseWriter, r *http.Request) {
	if isHXRequest(r) {
		h.render.HTMX(w, http.StatusOK, "profile.gohtml", nil)
		return
	}

	ctx := r.Context()
	u, _ := anor.UserFromContext(ctx)
	h.render.HTML(w, http.StatusOK, "profile.gohtml", ProfileViewData{User: u})
}
