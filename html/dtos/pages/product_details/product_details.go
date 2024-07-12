package product_details

import (
	"github.com/aliml92/anor/html/dtos/pages/product_details/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/html/dtos/shared"
)

type Base struct {
	partials.Header
	Content
}

type Content struct {
	shared.CategoryBreadcrumb
	components.ProductMain
	components.ProductSpecs
	ProductVariantMatrix interface{}
	// related products
}
