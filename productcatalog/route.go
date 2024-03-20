package productcatalog

import (
	"net/http"
)

func RegisterRoutes(h *Handler, mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.authInjector(h.NotFound(h.HomeView)))
	mux.HandleFunc("GET /products/{slug}", h.authInjector(h.ProductView))
	mux.HandleFunc("GET /categories/{slug}", h.authInjector(h.CategoryView))
	mux.HandleFunc("POST /search", h.SearchQuerySuggestions)
}
