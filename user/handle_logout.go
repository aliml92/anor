package user

import (
	"net/http"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.session.Clear(ctx)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	h.redirect(w, "/")
}
