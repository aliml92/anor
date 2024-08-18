package product

import (
	"github.com/aliml92/anor"
)

// RegisterRoutes sets up the routing for the product-related endpoints.
// It maps HTTP methods and paths to their respective handler functions.
func RegisterRoutes(h *Handler, router *anor.Router) {
	router.HandleFunc("POST /analytics/pl/views", h.CountViewOnProductListingView)
	router.HandleFunc("GET /products/{handle}", h.ProductDetailsView)
	router.HandleFunc("GET /categories/{handle}", h.ProductListingsView)
	router.HandleFunc("GET /search", h.SearchListingsView)
	router.HandleFunc("GET /search-query-suggestions", h.SearchQuerySuggestionsView)
	router.HandleFunc("GET /trending-products", h.TrendingProductsView)
	router.HandleFunc("GET /", h.HomeView, h.notFoundMiddleware)

	// sample handler to check the templates
	//router.HandleFunc("GET /checkout/address", h.CheckoutAddress)
	//router.HandleFunc("GET /checkout/address-add", h.CheckoutAddressUpdate)
	//router.HandleFunc("GET /checkout/payment", h.CheckoutPaymentMethod)
	//router.HandleFunc("GET /checkout/confirmnew", h.CheckoutConfirmNew)
	//router.HandleFunc("GET /checkout/done", h.CheckoutDone)
}

//func (h *Handler) CheckoutAddress(w http.ResponseWriter, r *http.Request) {
//	h.view.Render(w, "pages/checkout/address/base.gohtml", nil)
//}
//
//func (h *Handler) CheckoutAddressUpdate(w http.ResponseWriter, r *http.Request) {
//	h.view.Render(w, "pages/checkout/address_add/base.gohtml", nil)
//}
//
//func (h *Handler) CheckoutPaymentMethod(w http.ResponseWriter, r *http.Request) {
//	h.view.Render(w, "pages/checkout/payment_method/base.gohtml", nil)
//}
//
//func (h *Handler) CheckoutConfirmNew(w http.ResponseWriter, r *http.Request) {
//	h.view.Render(w, "pages/checkout/confirmnew/base.gohtml", nil)
//}
//
//func (h *Handler) CheckoutDone(w http.ResponseWriter, r *http.Request) {
//	h.view.Render(w, "pages/checkout/redirectnew/base.gohtml", nil)
//}
