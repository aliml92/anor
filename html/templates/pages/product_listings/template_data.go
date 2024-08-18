package product_listings

import (
	"github.com/aliml92/anor/html/templates/pages/product_listings/components"
	"github.com/aliml92/anor/html/templates/shared"
)

type Content struct {
	CategoryBreadcrumb  shared.CategoryBreadcrumb
	SideCategoryList    components.SideCategoryList
	SidePriceRange      components.SidePriceRange
	SideBrandsCheckbox  components.SideFilterCheckbox
	SideRatingCheckbox  components.SideRatingCheckbox
	ProductListingsInfo components.ProductListingsInfo
	ProductGrid         shared.ProductGrid
	Pagination          components.Pagination
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
