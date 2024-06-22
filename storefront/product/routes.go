package product

import (
	"github.com/aliml92/anor"
)

func RegisterRoutes(h *Handler, router *anor.Router) {
	router.HandleFunc("POST /analytics/pl/views", h.CountViewOnProductListingView)
	router.HandleFunc("GET /products/{handle}", h.ProductDetailsView)
	router.HandleFunc("GET /categories/{handle}", h.ProductListingView)
	router.HandleFunc("GET /search", h.SearchView)
	router.HandleFunc("GET /search-query-suggestions", h.SearchQuerySuggestionsView)
	router.HandleFunc("GET /trending-products", h.TrendingProductsView)
	router.HandleFunc("GET /", h.HomeView, h.notFoundMiddleware)
}
