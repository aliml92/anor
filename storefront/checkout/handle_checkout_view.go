package checkout

import (
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/config"
	"github.com/shopspring/decimal"
	"net/url"
)

//func (h *Handler) CheckoutView(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	u := session.UserFromContext(ctx)
//
//	c, err := h.cartService.Get(ctx, u.CartID, true)
//	if err != nil {
//		if errors.Is(err, anor.ErrCartNotFound) {
//			//TODO: redirect to /cart with a flash message
//			http.Redirect(w, r, "/cart", http.StatusFound)
//			return
//		}
//		h.serverInternalError(w, err)
//		return
//	}
//
//	if len(c.CartItems) == 0 {
//		//TODO: redirect to /cart with a flash message
//		http.Redirect(w, r, "/cart", http.StatusFound)
//		return
//	}
//
//	// get order by cart id
//	// check if its status Pending and Pending
//	// if yes check its shipping_address_id
//	//     if shipping_address_id is null
//	//         then return /checkout/address
//	//     else
//
//	// else
//	//     then create order and redirect to address
//
//	if u.ShippingAddressID == 0 || u.BillingAddressID == 0 {
//		redirectURL := buildRedirectURL(h.cfg.Server, "/checkout/address")
//		w.Header().Set("HX-Redirect", redirectURL)
//		w.WriteHeader(http.StatusOK)
//		return
//	}
//
//	if u.StripeConfirmationTokenID == "" {
//		http.Redirect(w, r, "/checkout/payment-method", http.StatusSeeOther)
//		return
//	}
//
//}

// calculateTotalAmount calculates total price of cart, it is assumed that
// all cart items have same currency code
func calculateTotalAmount(cartItems []*anor.CartItem) decimal.Decimal {
	var total decimal.Decimal
	for _, cartItem := range cartItems {
		total = total.Add(cartItem.Price)
	}
	return total
}

// getCurrency makes a naive assumption about cart currency code, it is assumed that
// all cart items have same currency code
func getCurrency(cartItems []*anor.CartItem) string {
	if len(cartItems) > 0 {
		return cartItems[0].Currency
	}

	return "USD"
}

func calculateCartTotalInCents(cartItems []*anor.CartItem) int64 {
	var totalPrice decimal.Decimal
	for _, cartItem := range cartItems {
		totalPrice = totalPrice.Add(cartItem.Price)
	}

	cents := totalPrice.Mul(decimal.NewFromInt(100))

	return cents.IntPart()
}

func buildRedirectURL(cfg config.ServerConfig, path string) string {
	redirectURL := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   path,
	}

	if cfg.IsHTTPS {
		redirectURL.Scheme = "https"
	}

	return redirectURL.String()
}
