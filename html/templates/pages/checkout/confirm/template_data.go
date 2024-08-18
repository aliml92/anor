package confirm

import (
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/shared"
	"github.com/shopspring/decimal"
)

type Content struct {
	Stepper               shared.Stepper
	EstimatedDeliveryDate string
	ShippingCost          decimal.Decimal
	CartItems             []*anor.CartItem
	ShippingAddress       anor.Address
	BillingAddress        anor.Address
	SelectedPaymentMethod PaymentMethodSummary
	CartTotal             decimal.Decimal
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}

type PaymentMethodSummary struct {
	Type  string // e.g., "Credit Card", "PayPal", etc.
	Last4 string // Last 4 digits if applicable
}
