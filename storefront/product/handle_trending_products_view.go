package product

import (
	"encoding/json"
	"net/http"
)

// TODO: remove hardcoded trending products, implement it instead if necessary
var trendingProducts = []string{"swimming pools", "sunglasses", "camping gear", "portable fan"}

// TrendingProductsView handles requests for trending product suggestions on search dropdown list
func (h *Handler) TrendingProductsView(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(trendingProducts)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}
}
