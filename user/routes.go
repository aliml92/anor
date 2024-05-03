package user

import (
	"net/http"
)

func RegisterRoutes(h *Handler, mux *http.ServeMux) {
	mux.HandleFunc("GET /user/profile", h.authInjector(h.ProfileView))
	mux.HandleFunc("DELETE /user/logout", h.authInjector(h.Logout))
}
