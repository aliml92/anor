package checkout

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/pages/checkout/address/components"
	"github.com/aliml92/anor/session"
	"github.com/invopop/validation"
	"net/http"
	"strconv"
)

type SetOrderAddressForm struct {
	AddressID   int64
	AddressType string
}

func (f *SetOrderAddressForm) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	id := r.PostForm.Get("address-id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	f.AddressID = idInt
	f.AddressType = r.PostForm.Get("address-type")

	return nil
}

func (f *SetOrderAddressForm) Validate() error {
	err := validation.Errors{
		"address-id":   validation.Validate(f.AddressID, validation.Required),
		"address-type": validation.Validate(f.AddressType, validation.Required, validation.In("Shipping", "Billing")),
	}.Filter()

	return err
}

func (h *Handler) SetOrderAddress(w http.ResponseWriter, r *http.Request) {
	var f SetOrderAddressForm
	err := anor.BindValid(r, &f)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	a, err := h.addressService.Get(ctx, f.AddressID)
	if err != nil {
		if errors.Is(err, anor.ErrAddressNotFound) {
			err = fmt.Errorf("%w. It seems the requested address has been deleted", anor.ErrAddressNotFound)
			h.clientError(w, err, http.StatusNotFound)
			return
		}
		h.serverInternalError(w, err)
		return
	}

	if f.AddressType == "Shipping" {
		h.session.Put(ctx, session.ShippingAddressIDKey, f.AddressID)
	} else {
		h.session.Put(ctx, session.BillingAddressIDKey, f.AddressID)
	}

	c := components.ShippingInfo{
		ShippingAddress: a,
	}

	h.Render(w, r, "pages/checkout/address/components/shipping_info.gohtml", c)
}
