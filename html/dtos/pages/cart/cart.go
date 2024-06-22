package cart

import (
	"github.com/aliml92/anor/html/dtos/pages/cart/components"
	"github.com/aliml92/anor/html/dtos/partials"
)

type Base struct {
	partials.Header
	Content
}

type Content struct {
	components.CartItems
	components.CartSummary
}
