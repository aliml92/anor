package auth

import (
	"net/http"
)

func RegisterRoutes(h *Handler, mux *http.ServeMux) {
	mux.HandleFunc("GET /signup", h.SignupView)
	mux.HandleFunc("POST /signup", h.Signup)
	mux.HandleFunc("POST /signup/confirm", h.SignupConfirm)

	mux.HandleFunc("GET /signin", h.SigninView)
	mux.HandleFunc("POST /signin", h.Signin)
}
