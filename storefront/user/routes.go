package user

import (
	"github.com/aliml92/anor"
)

// RegisterRoutes sets up the routing for the user-related endpoints.
// It maps HTTP methods and paths to their respective handler functions.
func RegisterRoutes(h *Handler, router *anor.Router) {
	router.HandleFunc("GET /user", h.ProfileView)
	router.HandleFunc("GET /user/orders", h.OrdersView)
	router.HandleFunc("POST /user/addresses", h.AddAddress)
	router.HandleFunc("DELETE /user/logout", h.Logout)
}
