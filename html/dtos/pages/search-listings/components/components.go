package components

import "github.com/aliml92/anor"

type CategoryBreadcrumb struct {
	Category           anor.Category
	AncestorCategories []anor.Category
	RootCategories     []anor.Category
}

type SearchListingsInfo struct {
	Q               string
	TotalCategories int
	TotalProducts   int
	SortParam       anor.SortParam
	FilterParam     anor.FilterParam
}

type SideCategoryList struct {
	Category           anor.Category
	AncestorCategories []anor.Category
	ChildrenCategories []anor.Category
	SiblingCategories  []anor.Category
	RootCategories     []anor.Category
}

type SideFilterCheckbox struct {
	Q             string
	FilterName    string
	FilteringData anor.FilterParam
	FilterParam   anor.FilterParam
}

type SidePriceRange struct {
	Q             string
	FilterParam   anor.FilterParam
	FilteringData anor.FilterParam
}

type SideRatingCheckbox struct {
}

type Pagination struct {
	Q             string
	TotalPages    int
	CurrentPage   int
	TotalProducts int
}
