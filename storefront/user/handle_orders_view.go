package user

import (
	"github.com/aliml92/anor"
	ordersdata "github.com/aliml92/anor/html/templates/pages/orders"
	"github.com/aliml92/anor/relation"
	"github.com/aliml92/anor/session"
	"net/http"
)

func (h *Handler) OrdersView(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")

	ctx := r.Context()
	u := session.UserFromContext(ctx)

	var (
		orders []anor.Order
		err    error
	)

	withRelations := relation.New(relation.ShippingAddress, relation.StripeCardPayment, relation.OrderItems)
	orderListParams := anor.OrderListParams{
		UserID:        u.ID,
		WithRelations: withRelations,
		Page:          1,
		PageSize:      10,
	}

	switch filter {
	case "active":
		orders, err = h.orderService.ListActive(ctx, orderListParams)
	case "unpaid":
		orders, err = h.orderService.ListUnpaid(ctx, orderListParams)
	default:
		orders, err = h.orderService.List(ctx, orderListParams)
	}
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	c := ordersdata.Content{Orders: orders}
	h.Render(w, r, "pages/orders/content.gohtml", c)
}
