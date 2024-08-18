package success

import (
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/shared"
	"github.com/shopspring/decimal"
	"time"
)

type Content struct {
	Stepper               shared.Stepper
	OrderID               int64
	OrderCreateAt         time.Time
	OrderTotal            decimal.Decimal
	ShippingAddress       anor.Address
	EstimatedDeliveryDate string
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
