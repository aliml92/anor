package user

import (
	"github.com/aliml92/anor"
)

func RegisterRoutes(h *Handler, router *anor.Router) {
	router.HandleFunc("GET /user/profile", h.ProfileView)
	router.HandleFunc("DELETE /user/logout", h.Logout)
}
