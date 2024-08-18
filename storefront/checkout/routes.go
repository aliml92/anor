package checkout

import (
	"github.com/aliml92/anor"
)

// RegisterRoutes sets up the routing for the checkout-related endpoints.
// It maps HTTP methods and paths to their respective handler functions.
func RegisterRoutes(h *Handler, router *anor.Router) {
	//router.HandleFunc("GET /checkout", h.CheckoutView)
	router.HandleFunc("GET /checkout/address", h.AddressView)
	router.HandleFunc("POST /checkout/address", h.SetOrderAddress)
	router.HandleFunc("GET /checkout/payment-method", h.PaymentMethodView)
	router.HandleFunc("GET /checkout/confirm", h.ConfirmView)
	router.HandleFunc("GET /checkout/success", h.SuccessView)

	// json apis
	router.HandleFunc("GET /checkout/order-summary", h.RetrieveOrderSummary)
	router.HandleFunc("POST /checkout/stripe-ctoken", h.SaveStripeConfirmationToken)
	router.HandleFunc("POST /checkout/create-payment-intent", h.CreatePaymentIntent)
}
