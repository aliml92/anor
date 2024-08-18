package checkout

import (
	"github.com/aliml92/anor/html/templates/pages/checkout/success"
	"github.com/aliml92/anor/html/templates/shared"
	"net/http"
	"strconv"
)

func (h *Handler) SuccessView(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(r.URL.Query().Get("order_id"), 10, 64)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	o, err := h.orderService.Get(ctx, orderID, false)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	p, err := h.stripePaymentService.GetByOrderID(ctx, orderID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	sh, err := h.addressService.Get(ctx, o.ShippingAddressID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	c := success.Content{
		Stepper:               shared.Stepper{CurrentStep: 5},
		OrderID:               o.ID,
		OrderCreateAt:         o.CreatedAt,
		OrderTotal:            p.Amount,
		ShippingAddress:       sh,
		EstimatedDeliveryDate: generateFakeEstimatedDeliveryTime(),
	}

	h.Render(w, r, "pages/checkout/success/content.gohtml", c)
}
