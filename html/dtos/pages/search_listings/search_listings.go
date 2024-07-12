package search_listings

import (
	"github.com/aliml92/anor/html/dtos/pages/search_listings/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/html/dtos/shared"
)

type Base struct {
	partials.Header
	Content
}

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
