package cart

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alexedwards/scs/v2"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/dtos/partials"
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
	req := &AddToCartRequest{}

	err := anor.BindValid(r, req)
	if err != nil {
		if errors.Is(err, errInternal) {
			h.serverInternalError(w, err)
			return
		}
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u := anor.UserFromContext(ctx)

	var cartID int64

	if u != nil {
		cartID, err = h.getOrCreateCart(ctx, u.ID)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}
	} else {
		guestSession := h.session.Guest
		cartID, err = h.getOrCreateGuestCart(ctx, guestSession)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}
	}

	_, err = h.cartSvc.AddCartItem(ctx, cartID, anor.AddCartItemParam{
		VariantID: req.VariantID,
		Qty:       req.Qty,
	})
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	itemCount, err := h.cartSvc.CountCartItems(ctx, cartID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	w.Header().Add("HX-Trigger-After-Settle", `{"anor:showToast":"item added to your cart successfully"}`)
	v := partials.CartNavItem{
		HxSwapOOB:      "true",
		CartItemsCount: int(itemCount),
	}
	h.view.RenderComponent(w, "partials/header/cart_nav_item.gohtml", v)
}

func (h *Handler) getOrCreateCart(ctx context.Context, userID int64) (int64, error) {
	c, err := h.cartSvc.GetCart(ctx, userID, false)
	if errors.Is(err, anor.ErrNotFound) {
		cart, err := h.cartSvc.CreateCart(ctx, userID)
		if err != nil {
			return 0, err
		}

		return cart.ID, nil

	} else if err != nil {
		return 0, err
	}

	return c.ID, nil
}

func (h *Handler) getOrCreateGuestCart(ctx context.Context, guestSession *scs.SessionManager) (int64, error) {
	cartID := guestSession.GetInt64(ctx, "guest_cart_id")
	if cartID == 0 {
		cart, err := h.cartSvc.CreateCart(ctx, -1)
		if err != nil {
			return 0, err
		}
		guestSession.Put(ctx, "guest_cart_id", cart.ID)

		return cart.ID, nil
	}

	return cartID, nil
}
