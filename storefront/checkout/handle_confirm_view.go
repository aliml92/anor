package checkout

import (
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/pages/checkout/confirm"
	"github.com/aliml92/anor/html/templates/shared"
	"github.com/aliml92/anor/session"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/confirmationtoken"
	"math/rand"
	"net/http"
	"time"
)

func (h *Handler) ConfirmView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := session.UserFromContext(ctx)

	if u.CartID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if u.ShippingAddressID == 0 || u.BillingAddressID == 0 {
		http.Redirect(w, r, "/checkout/address", http.StatusSeeOther)
		return
	}

	if u.StripeConfirmationTokenID == "" {
		http.Redirect(w, r, "/checkout/payment-method", http.StatusSeeOther)
		return
	}

	c, err := h.cartService.Get(ctx, u.CartID, true)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	shippingAddr, err := h.addressService.Get(ctx, u.ShippingAddressID)
	if err != nil {
		h.serverInternalError(w, err)
	}

	var billingAddr anor.Address
	if u.ShippingAddressID != u.BillingAddressID {
		billingAddr, err = h.addressService.Get(ctx, u.BillingAddressID)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}
	} else {
		billingAddr = shippingAddr
	}

	stripe.Key = h.cfg.Stripe.SecretKey
	confirmationToken, err := confirmationtoken.Get(u.StripeConfirmationTokenID, &stripe.ConfirmationTokenParams{})
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	content := confirm.Content{
		Stepper: shared.Stepper{
			CurrentStep: 4,
		},
		EstimatedDeliveryDate: generateFakeEstimatedDeliveryTime(),
		ShippingCost:          decimal.Zero, //  TODO: add actual data later
		CartItems:             c.CartItems,
		ShippingAddress:       shippingAddr,
		BillingAddress:        billingAddr,
		SelectedPaymentMethod: confirm.PaymentMethodSummary{
			Type:  confirmationToken.PaymentMethodPreview.Card.Brand,
			Last4: confirmationToken.PaymentMethodPreview.Card.Last4,
		},
		CartTotal: calculateTotalAmount(c.CartItems),
	}

	h.Render(w, r, "pages/checkout/confirm/content.gohtml", content)
}

func generateFakeEstimatedDeliveryTime() string {
	now := time.Now()

	minDays := 3
	startDays := rand.Intn(3) + minDays

	windowDays := rand.Intn(3) + 2

	startDate := now.AddDate(0, 0, startDays)
	endDate := startDate.AddDate(0, 0, windowDays)

	return fmt.Sprintf("%s - %s",
		startDate.Format("Mon, 02/01"),
		endDate.Format("Mon, 02/01"))
}
