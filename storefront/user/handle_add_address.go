package user

import (
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/config"
	"github.com/aliml92/anor/session"
	"github.com/invopop/validation"
	"net/http"
	"net/url"
)

type AddAddressForm struct {
	Name          string
	AddressLine1  string
	AddressLine2  string
	City          string
	StateProvince string
	PostalCode    string
	Country       string
	DefaultFor    string
	SelectedAs    string
}

func (f *AddAddressForm) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	f.Name = r.PostForm.Get("name")
	f.AddressLine1 = r.PostForm.Get("address-line-1")
	f.AddressLine2 = r.PostForm.Get("address-line-2")
	f.City = r.PostForm.Get("city")
	f.StateProvince = r.PostForm.Get("state-province")
	f.PostalCode = r.PostForm.Get("postal-code")
	f.Country = r.PostForm.Get("country")
	f.DefaultFor = r.PostForm.Get("default-for")
	f.SelectedAs = r.PostForm.Get("selected-as")

	return nil
}

func (f *AddAddressForm) Validate() error {
	return validation.Errors{
		"name":           validation.Validate(f.Name, validation.Required, validation.Length(1, 100)),
		"address-line-1": validation.Validate(f.AddressLine1, validation.Required, validation.Length(1, 100)),
		"address-line-2": validation.Validate(f.AddressLine2, validation.Length(0, 100)),
		"city":           validation.Validate(f.City, validation.Required, validation.Length(1, 50)),
		"state-province": validation.Validate(f.StateProvince, validation.Required, validation.Length(1, 50)),
		"postal-code":    validation.Validate(f.PostalCode, validation.Required, validation.Length(1, 20)),
		"country":        validation.Validate(f.Country, validation.Required, validation.Length(1, 50)),
		"default-for":    validation.Validate(f.DefaultFor, validation.In("Shipping", "Billing", "")),
		"selected-as":    validation.Validate(f.SelectedAs, validation.In("Shipping", "Billing", "Both", "")),
	}.Filter()
}

func (h *Handler) AddAddress(w http.ResponseWriter, r *http.Request) {
	var f AddAddressForm
	err := anor.BindValid(r, &f)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u := session.UserFromContext(ctx)

	var defaultFor anor.AddressDefaultType
	switch f.SelectedAs {
	case "Shipping":
		defaultFor = anor.AddressDefaultTypeShipping
	case "Billing":
		defaultFor = anor.AddressDefaultTypeBilling
	case "Both":
	default:
		defaultFor = anor.AddressDefaultType(f.DefaultFor)
	}

	a, err := h.addressService.Create(ctx, anor.AddressCreateParams{
		UserID:        u.ID,
		DefaultFor:    defaultFor,
		Name:          f.Name,
		AddressLine1:  f.AddressLine1,
		AddressLine2:  f.AddressLine2,
		City:          f.City,
		StateProvince: f.StateProvince,
		PostalCode:    f.PostalCode,
		Country:       f.Country,
	})
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	switch f.SelectedAs {
	case "Shipping":
		h.session.Put(ctx, session.ShippingAddressIDKey, a.ID)
	case "Billing":
		h.session.Put(ctx, session.BillingAddressIDKey, a.ID)
	case "Both":
		h.session.Put(ctx, session.ShippingAddressIDKey, a.ID)
		h.session.Put(ctx, session.BillingAddressIDKey, a.ID)
	default:
	}

	redirectURL := r.URL.Query().Get("redirect_url")
	if redirectURL != "" {
		redirectURL = buildRedirectURL(h.cfg.Server, redirectURL)
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
		return
	}
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
