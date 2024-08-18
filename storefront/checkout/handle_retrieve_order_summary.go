package checkout

import (
	"encoding/json"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/session"
	"github.com/shopspring/decimal"
	"net/http"
	"strings"
)

func (h *Handler) RetrieveOrderSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := session.UserFromContext(ctx)

	c, err := h.cartService.Get(ctx, u.CartID, true)
	if err != nil {
		h.jsonServerInternalError(w, err)
		return
	}

	totalAmount := calculateTotalAmount(c.CartItems)
	amountInCents := totalAmount.Mul(decimal.NewFromInt(100)).IntPart()

	response := struct {
		AmountInCents int64  `json:"amount"`
		Currency      string `json:"currency"`
	}{
		AmountInCents: amountInCents,
		Currency:      strings.ToLower(anor.DefaultCurrency),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.jsonServerInternalError(w, err)
		return
	}
}
