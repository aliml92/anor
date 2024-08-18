package address

import (
	"github.com/aliml92/anor/html/templates/pages/checkout/address/components"
	"github.com/aliml92/anor/html/templates/shared"
)

type Content struct {
	Stepper                     shared.Stepper
	ShippingAddressKindSelector components.ShippingAddressKindSelector
	ShippingInfo                components.ShippingInfo
	BillingInfo                 components.BillingInfo
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
