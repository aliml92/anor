package user

import (
	"net/http"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.session.Auth.Destroy(ctx); err != nil {
		h.serverInternalError(w, err)
		return
	}

	h.redirect(w, "/")
}
