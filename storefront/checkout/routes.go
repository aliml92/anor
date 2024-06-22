package checkout

import (
	"github.com/aliml92/anor"
)

func RegisterRoutes(h *Handler, router *anor.Router) {
	router.HandleFunc("GET /checkout", h.GetCheckoutView)
	router.HandleFunc("POST /checkout/get-pi-client-secret", h.GetPIClientSecret)
	router.HandleFunc("GET /checkout/redirect", h.GetCheckoutRedirectView)
}
