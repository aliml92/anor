package components

import "github.com/aliml92/anor"

type ShippingInfo struct {
	ShippingAddress anor.Address
}

func (ShippingInfo) GetTemplateFilename() string {
	return "shipping_info.gohtml"
}

type BillingInfo struct {
	SameAsShipping bool
	BillingAddress anor.Address
}

func (BillingInfo) GetTemplateFilename() string {
	return "billing_info.gohtml"
}

type ShippingAddressForm struct {
	ShippingAddressKindSelector
}

func (ShippingAddressForm) GetTemplateFilename() string {
	return "shipping_address_form.gohtml"
}

type ShippingAddressKindSelector struct {
	NoUserAddresses     bool
	SelectedAddressKind string
}

func (ShippingAddressKindSelector) GetTemplateFilename() string {
	return "shipping_address_kind_selector.gohtml"
}
