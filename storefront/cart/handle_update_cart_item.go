package cart

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/dtos/pages/cart/components"
	"github.com/invopop/validation"
	"net/http"
	"strconv"
)

type UpdateCartItemRequest struct {
	CartItemID int64
	Qty        int
}

func (req *UpdateCartItemRequest) Bind(r *http.Request) error {
	cartItemIDStr := r.PathValue("id")
	cartItemID, err := strconv.ParseInt(cartItemIDStr, 10, 64)
	if err != nil {
		return err
	}
	req.CartItemID = cartItemID

	qty := r.FormValue("qty")
	fmt.Println("qty:", qty)
	req.Qty, err = strconv.Atoi(qty)
	if err != nil {
		return err
	}

	return nil
}

func (req *UpdateCartItemRequest) Validate() error {
	err := validation.Errors{
		"id":  validation.Validate(req.CartItemID, validation.Required, validation.Min(1)),
		"qty": validation.Validate(req.Qty, validation.Required, validation.Min(1)),
	}.Filter()

	return err
}

func (h *Handler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	req := &UpdateCartItemRequest{}

	err := bindValid(r, req)
	if err != nil {
		if errors.Is(err, errInternal) {
			h.serverInternalError(w, err)
			return
		}
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.cartSvc.UpdateCartItem(ctx, req.CartItemID, anor.UpdateCartItemParam{Qty: req.Qty})
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	u := anor.UserFromContext(ctx)
	c, err := h.cartSvc.GetCart(ctx, u.ID, true)
	if err != nil && !errors.Is(err, anor.ErrNotFound) {
		h.serverInternalError(w, err)
		return
	}

	w.Header().Add("HX-Trigger-After-Settle", `{"anor:showToast":"item qty updated successfully"}`)
	v := components.CartSummary{
		HxSwapOOB:      r.Header.Get("Hx-Swap-OOB"),
		CartItemsCount: len(c.CartItems),
		TotalAmount:    getTotalPrice(c.CartItems),
		CurrencyCode:   getCurrency(c.CartItems),
	}

	h.view.RenderComponent(w, "pages/cart/components/cart-summary.gohtml", v)
}
