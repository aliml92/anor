package auth

import (
	"net/http"
)

func RegisterRoutes(h *Handler, mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/signup", h.SignupView)
	mux.HandleFunc("POST /auth/signup", h.Signup)

	mux.HandleFunc("POST /auth/confirmation", h.SignupConfirmation)
	mux.HandleFunc("POST /auth/confirmation/resend", h.ResendOTP)

	mux.HandleFunc("GET /auth/signin", h.SigninView)
	mux.HandleFunc("POST /auth/signin", h.Signin)

	mux.HandleFunc("GET /auth/forgot-password", h.ForgotPasswordView)
	mux.HandleFunc("POST /auth/send-reset-link", h.SendResetPasswordLink)

	mux.HandleFunc("GET /auth/reset-password", h.ResetPasswordView)
	mux.HandleFunc("POST /auth/reset-password", h.ResetPassword)
}
