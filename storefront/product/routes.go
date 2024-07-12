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
}
