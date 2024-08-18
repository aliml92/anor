package checkout

import "C"
import (
	"context"
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/pages/checkout/address"
	"github.com/aliml92/anor/html/templates/pages/checkout/address/components"
	"github.com/aliml92/anor/html/templates/shared"
	"github.com/aliml92/anor/session"
	"net/http"
)

func (h *Handler) AddressView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	targetTemplate := r.Header.Get("X-Target-Template")

	switch targetTemplate {
	case "empty-address-form":
		//c := address_add.Content{
		//	IsEmpty: true,
		//}
		//h.view.Render(w, "pages/checkout/address_add/content.gohtml", c)
		//return

	case "address-form":
		//u := session.UserFromContext(ctx)
		//selectedShipping, err := h.getAddressOrDefault(ctx, u.ShippingAddressID, u.ID, anor.AddressDefaultTypeShipping)
		//if err != nil {
		//	h.serverInternalError(w, err)
		//	return
		//}
		//
		//shippingAddresses, err := h.listAddresses(ctx, u.ID)
		//if err != nil {
		//	h.serverInternalError(w, err)
		//	return
		//}

		//c := address_add.Content{
		//	IsEmpty:           false,
		//	ShippingAddresses: shippingAddresses,
		//}
		//if !selectedShipping.IsZero() {
		//	c.SelectedShippingAddress = selectedShipping
		//}

		//h.view.Render(w, "pages/checkout/address_add/content.gohtml", c)
		//return

	default:
		u := session.UserFromContext(ctx)

		c := address.Content{
			Stepper:                     shared.Stepper{CurrentStep: 2},
			ShippingAddressKindSelector: components.ShippingAddressKindSelector{},
			BillingInfo:                 components.BillingInfo{},
			ShippingInfo:                components.ShippingInfo{},
		}

		shipping, err := h.getAddressOrDefault(ctx, u.ShippingAddressID, u.ID, anor.AddressDefaultTypeShipping)
		billing, err := h.getAddressOrDefault(ctx, u.BillingAddressID, u.ID, anor.AddressDefaultTypeBilling)
		if errors.Is(err, anor.ErrAddressNotFound) {
			addresses, err := h.listAddresses(ctx, u.ID)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}

			if len(addresses) > 0 {
				if shipping.IsZero() {
					shipping = addresses[0]
					h.session.Put(ctx, session.ShippingAddressIDKey, shipping.ID)
				}
				if billing.IsZero() {
					billing = addresses[0]
					h.session.Put(ctx, session.BillingAddressIDKey, billing.ID)
				}
			} else {
				c.ShippingAddressKindSelector.NoUserAddresses = true
			}
		} else if err != nil {
			h.serverInternalError(w, err)
			return
		}

		if !shipping.IsZero() || billing.IsZero() {
			c.BillingInfo.SameAsShipping = billing.Equals(shipping) == true
		}
		c.ShippingInfo.ShippingAddress = shipping
		c.BillingInfo.BillingAddress = billing

		h.Render(w, r, "pages/checkout/address/content.gohtml", c)
	}

}

func (h *Handler) getAddressOrDefault(ctx context.Context, addressID int64, userID int64, defaultFor anor.AddressDefaultType) (anor.Address, error) {
	if addressID == 0 {
		return h.addressService.GetDefault(ctx, userID, defaultFor)
	} else {
		addr, err := h.addressService.Get(ctx, addressID)
		if errors.Is(err, anor.ErrAddressNotFound) {
			addr, err = h.addressService.GetDefault(ctx, userID, defaultFor)
		}
		return addr, err
	}
}

func (h *Handler) listAddresses(ctx context.Context, userID int64) ([]anor.Address, error) {
	return h.addressService.List(ctx, anor.AddressListParams{
		UserID:   userID,
		Page:     1,
		PageSize: 5,
	})
}
