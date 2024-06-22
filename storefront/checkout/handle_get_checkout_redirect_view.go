package checkout

import (
	"github.com/aliml92/anor"
	checkout_redirect "github.com/aliml92/anor/html/dtos/pages/checkout-redirect"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/invopop/validation"
	"log"
	"net/http"
	"net/url"
)

type RedirectQuery struct {
	PaymentIntent             string
	PaymentIntentClientSecret string
	RedirectStatus            string
}

func (q *RedirectQuery) Bind(r *http.Request) error {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return err
	}

	paymentIntent := values.Get("payment_intent")
	paymentIntentClientSecret := values.Get("payment_intent_client_secret")
	redirectStatus := values.Get("redirect_status")

	q.PaymentIntent = paymentIntent
	q.PaymentIntentClientSecret = paymentIntentClientSecret
	q.RedirectStatus = redirectStatus

	return nil
}

func (q *RedirectQuery) Validate() error {
	err := validation.Errors{
		"payment_intent":               validation.Validate(q.PaymentIntent, validation.Required),
		"payment_intent_client_secret": validation.Validate(q.PaymentIntentClientSecret, validation.Required),
		"redirect_status":              validation.Validate(q.RedirectStatus, validation.Required),
	}.Filter()

	return err
}

type GetCheckoutRedirectViewData struct {
	RedirectQuery
	User  *anor.User
	Order anor.Order
}

func (h *Handler) GetCheckoutRedirectView(w http.ResponseWriter, r *http.Request) {
	q := &RedirectQuery{}

	err := bindValid(r, q)
	if err != nil {
		h.logClientError(err)
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u := anor.UserFromContext(ctx)

	v := &GetCheckoutRedirectViewData{
		RedirectQuery: *q,
		User:          u,
	}

	var order anor.Order

	switch q.RedirectStatus {
	case "succeeded":
		var (
			cart anor.Cart
			err  error
		)

		if u != nil {
			cart, err = h.cartSvc.GetCart(ctx, u.ID, true)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}

			// create order
			newOrder, err := h.orderSvc.ConvertCartToOrder(ctx, cart, q.PaymentIntent)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}
			order = newOrder
			log.Printf("total amount: %v\n", v.Order.TotalAmount)
		} else {
			//TODO: Handle guest checkout differently
		}
	}

	crc := checkout_redirect.Content{
		Order: order,
	}

	headerContent := partials.Header{User: u}
	if u != nil {
		ac, err := h.userSvc.GetUserActivityCounts(ctx, u.ID)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

		headerContent.ActiveOrdersCount = ac.ActiveOrdersCount
		headerContent.WishlistItemsCount = ac.WishlistItemsCount
		headerContent.CartItemsCount = ac.CartItemsCount

	} else {
		cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
		if cartId != 0 {

			guestCartItemCount, err := h.cartSvc.CountCartItems(ctx, cartId)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}

			headerContent.CartItemsCount = int(guestCartItemCount)
		}
	}

	crb := checkout_redirect.Base{
		Header:  headerContent,
		Content: crc,
	}

	h.view.Render(w, "pages/checkout/base.gohtml", crb)
}
