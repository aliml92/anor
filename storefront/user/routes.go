package user

import (
	"github.com/aliml92/anor"
)

// RegisterRoutes sets up the routing for the user-related endpoints.
// It maps HTTP methods and paths to their respective handler functions.
func RegisterRoutes(h *Handler, router *anor.Router) {
	router.HandleFunc("GET /user/profile", h.ProfileView)
	router.HandleFunc("DELETE /user/logout", h.Logout)
}
