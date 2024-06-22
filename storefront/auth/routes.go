package auth

import (
	"github.com/aliml92/anor"
)

func RegisterRoutes(h *Handler, router *anor.Router) {
	authRouter := router.Group("/auth")
	authRouter.HandleFunc("GET /signup", h.SignupView, h.RedirectAuthUserMiddleware)
	authRouter.HandleFunc("POST /signup", h.Signup)

	authRouter.HandleFunc("POST /confirmation", h.SignupConfirmation)
	authRouter.HandleFunc("POST /confirmation/resend", h.ResendOTP)

	authRouter.HandleFunc("GET /signin", h.SigninView, h.RedirectAuthUserMiddleware)
	authRouter.HandleFunc("POST /signin", h.Signin)

	authRouter.HandleFunc("GET /forgot-password", h.ForgotPasswordView)
	authRouter.HandleFunc("POST /send-reset-link", h.SendResetPasswordLink)

	authRouter.HandleFunc("GET /reset-password", h.ResetPasswordView)
	authRouter.HandleFunc("POST /reset-password", h.ResetPassword)
}
