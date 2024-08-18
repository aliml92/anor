package product_details

import (
	"github.com/aliml92/anor/html/templates/pages/product_details/components"
	"github.com/aliml92/anor/html/templates/shared"
)

type Content struct {
	CategoryBreadcrumb   shared.CategoryBreadcrumb
	ProductMain          components.ProductMain
	ProductSpecs         components.ProductSpecs
	ProductVariantMatrix interface{}
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
