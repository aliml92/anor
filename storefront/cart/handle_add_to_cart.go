package cart

import (
	"encoding/json"
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/shared/header/components"
	"github.com/aliml92/anor/session"
	"github.com/invopop/validation"
	"io"
	"net/http"
)

var errInternal = errors.New("internal error")

type AddToCartRequest struct {
	VariantID int64 `json:"product_variant_id"`
	Qty       int   `json:"qty"`
}

func (req *AddToCartRequest) Bind(r *http.Request) error {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return errors.Join(err, errInternal)
	}

	err = json.Unmarshal(b, req)
	if err != nil {
		return err
	}

	return nil
}

func (req *AddToCartRequest) Validate() error {
	err := validation.Errors{
		"product_variant_id": validation.Validate(req.VariantID, validation.Required, validation.Min(1)),
		"qty":                validation.Validate(req.Qty, validation.Required, validation.Min(1)),
	}.Filter()

	return err
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req AddToCartRequest
	err := anor.BindValid(r, &req)
	if err != nil {
		if errors.Is(err, errInternal) {
			h.serverInternalError(w, err)
			return
		}
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u := session.UserFromContext(ctx)

	cartID := u.CartID
	if cartID == 0 {
		params := anor.CartCreateParams{}

		if u.IsAuth {
			params.UserID = u.ID
		} else {
			params.ExpiresAt = h.session.GetExpiry(ctx)
		}

		cart, err := h.cartService.Create(ctx, params)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

		cartID = cart.ID

		// save created cart id in the session
		h.session.Put(ctx, session.CartIDKey, cart.ID)
	}

	_, err = h.cartService.AddItem(ctx, anor.CartItemAddParams{
		CartID:    cartID,
		VariantID: req.VariantID,
		Qty:       req.Qty,
	})
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	itemCount, err := h.cartService.CountItems(ctx, cartID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	w.Header().Add("HX-Trigger-After-Settle", `{"anor:showToast":"item added to your cart successfully"}`)
	c := components.CartNavItem{
		HxSwapOOB:      true,
		CartItemsCount: int(itemCount),
	}

	h.Render(w, r, "shared/header/components/cart_nav_item.gohtml", c)
}
