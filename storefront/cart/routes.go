package cart

import (
	"github.com/aliml92/anor"
)

func RegisterRoutes(h *Handler, router *anor.Router) {
	router.HandleFunc("GET /cart", h.CartView)
	router.HandleFunc("POST /cart", h.AddToCart)
	router.HandleFunc("PATCH /cart/item/{id}", h.UpdateCartItem, h.VerifyCartItemOwnership)
	router.HandleFunc("DELETE /cart/item/{id}", h.RemoveCartItem, h.VerifyCartItemOwnership)
}
