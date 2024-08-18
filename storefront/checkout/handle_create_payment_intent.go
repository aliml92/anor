package checkout

import (
	"encoding/json"
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/session"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/confirmationtoken"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"net/http"
	"strconv"
	"time"
)

type CreatePaymentIntentResponse struct {
	ClientSecret string                     `json:"client_secret"`
	Status       stripe.PaymentIntentStatus `json:"status"`
	OrderID      int64                      `json:"order_id,omitempty"`
}

func (h *Handler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := session.UserFromContext(ctx)

	c, err := h.cartService.Get(ctx, u.CartID, true)
	if err != nil {
		if errors.Is(err, anor.ErrCartNotFound) {
			http.Redirect(w, r, "/cart", http.StatusFound)
			return
		}
		h.serverInternalError(w, err)
		return
	}

	if len(c.CartItems) == 0 {
		http.Redirect(w, r, "/cart", http.StatusFound)
		return
	}

	stripe.Key = h.cfg.Stripe.SecretKey

	cartTotal := calculateCartTotalInCents(c.CartItems)
	intentParams := &stripe.PaymentIntentParams{
		Confirm:           anor.Bool(true),
		Amount:            anor.Int64(cartTotal),
		Currency:          anor.String(anor.DefaultCurrency),
		ConfirmationToken: anor.String(u.StripeConfirmationTokenID),
		Metadata: map[string]string{
			"cart_id": strconv.FormatInt(c.ID, 10),
		},
		ReturnURL: anor.String("http://localhost:8008/checkout/redirect"),
	}

	intent, err := paymentintent.New(intentParams)
	if err != nil {
		h.jsonServerInternalError(w, err)
		return
	}

	var orderID int64
	if intent.Status == stripe.PaymentIntentStatusSucceeded {
		orderParams := anor.OrderCreateParams{
			Cart:              c,
			PaymentMethod:     anor.PaymentMethodStripeCard, // TODO: set actual data from intent's payment_method
			PaymentStatus:     anor.PaymentStatusPaid,
			ShippingAddressID: u.ShippingAddressID,
			IsPickup:          false,
			Amount:            calculateTotalAmount(c.CartItems),
			Currency:          anor.DefaultCurrency,
		}
		orderID, err = h.orderService.Create(ctx, orderParams)
		if err != nil {
			h.jsonServerInternalError(w, err)
			return
		}

		confirmationToken, err := confirmationtoken.Get(u.StripeConfirmationTokenID, &stripe.ConfirmationTokenParams{})
		if err != nil {
			h.jsonServerInternalError(w, err)
			return
		}

		var last4 string
		var cardBrand string
		if confirmationToken.PaymentMethodPreview.Card != nil {
			last4 = confirmationToken.PaymentMethodPreview.Card.Last4
			cardBrand = confirmationToken.PaymentMethodPreview.Card.DisplayBrand
		}
		err = h.stripePaymentService.Create(ctx, anor.StripePaymentCreateParams{
			OrderID:          orderID,
			UserID:           u.ID,
			BillingAddressID: u.BillingAddressID,
			PaymentIntentID:  intent.ID,
			PaymentMethodID:  intent.PaymentMethod.ID,
			Amount:           convertToDecimal(intent.Amount),
			Currency:         string(intent.Currency),
			Status:           string(intent.Status),
			ClientSecret:     intent.ClientSecret,
			LastError:        intent.LastPaymentError.Error(),
			CardLast4:        last4,
			CardBrand:        cardBrand,
			CreatedAt:        convertInt64ToTime(intent.Created),
		})
		if err != nil {
			h.jsonServerInternalError(w, err)
			return
		}

		err = h.cartService.Update(ctx, u.CartID, anor.CartUpdateParams{
			Status: anor.CartStatusCompleted,
		})
		if err != nil {
			h.jsonServerInternalError(w, err)
			return
		}

		h.session.Remove(ctx, session.CartIDKey)
		h.session.Remove(ctx, session.StripeConfirmationTokenIDKey)
	}

	response := CreatePaymentIntentResponse{
		ClientSecret: intent.ClientSecret,
		Status:       intent.Status,
		OrderID:      orderID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.jsonServerInternalError(w, err)
		return
	}

}

func convertToDecimal(amount int64) decimal.Decimal {
	// Convert int64 to Decimal
	decimalAmount := decimal.NewFromInt(amount)

	// Divide by 100
	result := decimalAmount.Div(decimal.NewFromInt(100))

	return result
}

func convertInt64ToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}
