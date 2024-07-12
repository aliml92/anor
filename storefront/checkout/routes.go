package checkout

import (
	"github.com/aliml92/anor"
)

// RegisterRoutes sets up the routing for the checkout-related endpoints.
// It maps HTTP methods and paths to their respective handler functions.
func RegisterRoutes(h *Handler, router *anor.Router) {
	router.HandleFunc("GET /checkout", h.GetCheckoutView)
	router.HandleFunc("POST /checkout/get-pi-client-secret", h.GetPIClientSecret)
	router.HandleFunc("GET /checkout/redirect", h.GetCheckoutRedirectView)
}
