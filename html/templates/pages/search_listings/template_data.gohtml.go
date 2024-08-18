package search_listings

import (
	"github.com/aliml92/anor/html/templates/pages/search_listings/components"
	"github.com/aliml92/anor/html/templates/shared"
)

type Content struct {
	CategoryBreadcrumb components.CategoryBreadcrumb
	SideCategoryList   components.SideCategoryList
	SidePriceRange     components.SidePriceRange
	SideBrandsCheckbox components.SideFilterCheckbox
	SideRatingCheckbox components.SideRatingCheckbox
	SearchListingsInfo components.SearchListingsInfo
	ProductGrid        shared.ProductGrid
	Pagination         components.Pagination
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
