package checkout

import (
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/dtos/pages/checkout"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"net/http"
	"strconv"
)

type GetCheckoutViewData struct {
	User *anor.User
	Cart *anor.Cart
}

func (h *Handler) GetCheckoutView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := anor.UserFromContext(ctx)

	var (
		cartItems []*anor.CartItem
		err       error
	)

	stripe.Key = h.cfg.Stripe.SecretKey
	if u != nil {
		c, err := h.cartSvc.GetCart(ctx, u.ID, true)
		if err != nil {
			if errors.Is(err, anor.ErrNotFound) {
				http.Redirect(w, r, "/carts", http.StatusFound)
				return
			}
			h.serverInternalError(w, err)
			return
		}

		if c.PIClientSecret == "" {
			// create payment intent
			paymentIntentParams := &stripe.PaymentIntentParams{
				Amount: stripe.Int64(calculateOrderAmount(c.CartItems)),
				// TODO: for now assume all products price in usd
				Currency: stripe.String(string(stripe.CurrencyUSD)),
				Metadata: map[string]string{
					"cart_id": strconv.FormatInt(c.ID, 10),
				},
			}

			pi, err := paymentintent.New(paymentIntentParams)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}
			c.PIClientSecret = pi.ClientSecret

			// update cart with payment intent's client secret
			if err := h.cartSvc.UpdateCart(ctx, c); err != nil {
				h.serverInternalError(w, err)
				return
			}
		}
		cartItems = c.CartItems
		//return cartItems
	} else {
		cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
		if cartId != 0 {
			cartItems, err = h.cartSvc.GetGuestCartItems(ctx, cartId)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}
		}
	}

	c := checkout.Content{Cart: anor.Cart{
		CartItems:    cartItems,
		TotalAmount:  getTotalPrice(cartItems),
		CurrencyCode: getCurrency(cartItems),
	},
	}

	h.Render(w, r, "pages/checkout", c)
}

// getTotalPrice calculates total price of cart, it is assumed that
// all cart items have same currency code
func getTotalPrice(cartItems []*anor.CartItem) decimal.Decimal {
	var totalPrice decimal.Decimal
	for _, cartItem := range cartItems {
		totalPrice = totalPrice.Add(cartItem.Price)
	}
	return totalPrice
}

// getCurrency makes a naive assumption about cart currency code, it is assumed that
// all cart items have same currency code
func getCurrency(cartItems []*anor.CartItem) string {
	if len(cartItems) > 0 {
		return cartItems[0].CurrencyCode
	}

	return "USD"
}

func calculateOrderAmount(cartItems []*anor.CartItem) int64 {
	var totalPrice decimal.Decimal
	for _, cartItem := range cartItems {
		totalPrice = totalPrice.Add(cartItem.Price)
	}

	cents := totalPrice.Mul(decimal.NewFromInt(100))

	return cents.IntPart()
}
