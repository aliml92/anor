package checkout

import (
	paymentmethod "github.com/aliml92/anor/html/templates/pages/checkout/payment_method"
	"github.com/aliml92/anor/html/templates/shared"
	"github.com/aliml92/anor/session"
	"net/http"
)

func (h *Handler) PaymentMethodView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := session.UserFromContext(ctx)

	if u.CartID == 0 {
		http.Redirect(w, r, "/cart", http.StatusFound)
		return
	}

	if u.ShippingAddressID == 0 || u.BillingAddressID == 0 {
		http.Redirect(w, r, "/checkout/address", http.StatusFound)
		return
	}

	c := paymentmethod.Content{
		Stepper: shared.Stepper{
			CurrentStep: 3,
		},
	}
	h.Render(w, r, "pages/checkout/payment_method/content.gohtml", c)
}
