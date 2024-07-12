package components

import (
	"github.com/aliml92/anor"
)

type Pagination struct {
	CategoryPath  string
	TotalPages    int
	CurrentPage   int
	TotalProducts int
}

type ProductListingsInfo struct {
	CategoryPath  string
	TotalProducts int
	SortParam     anor.SortParam
	FilterParam   anor.FilterParam
}

type SideCategoryList struct {
	Category           anor.Category
	AncestorCategories []anor.Category
	ChildrenCategories []anor.Category
	SiblingCategories  []anor.Category
}

type SideFilterCheckbox struct {
	FilterName    string
	CategoryPath  string
	FilteringData anor.FilterParam
	FilterParam   anor.FilterParam
}

type SidePriceRange struct {
	CategoryPath  string
	FilterParam   anor.FilterParam
	FilteringData anor.FilterParam
}

type SideRatingCheckbox struct {
}
