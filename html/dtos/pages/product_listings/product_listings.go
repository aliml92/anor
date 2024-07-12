package product_listings

import (
	"github.com/aliml92/anor/html/dtos/pages/product_listings/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/html/dtos/shared"
)

type Base struct {
	partials.Header
	Content
}

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
