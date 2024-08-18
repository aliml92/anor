package cart

import "github.com/aliml92/anor/html/templates/pages/cart/components"

type Content struct {
	CartItems   components.CartItems
	CartSummary components.CartSummary
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
