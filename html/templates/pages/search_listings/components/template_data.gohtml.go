package components

import "github.com/aliml92/anor"

type CategoryBreadcrumb struct {
	Category           anor.Category
	AncestorCategories []anor.Category
	RootCategories     []anor.Category
}

func (CategoryBreadcrumb) GetTemplateFilename() string {
	return "category_breadcrumb.gohtml"
}

type SearchListingsInfo struct {
	Q               string
	TotalCategories int
	TotalProducts   int
	SortParam       anor.SortParam
	FilterParam     anor.FilterParam
}

func (SearchListingsInfo) GetTemplateFilename() string {
	return "search_listings_info.gohtml"
}

type SideCategoryList struct {
	Category           anor.Category
	AncestorCategories []anor.Category
	ChildrenCategories []anor.Category
	SiblingCategories  []anor.Category
	RootCategories     []anor.Category
}

func (SideCategoryList) GetTemplateFilename() string {
	return "side_category_list.gohtml"
}

type SideFilterCheckbox struct {
	Q             string
	FilterName    string
	FilteringData anor.FilterParam
	FilterParam   anor.FilterParam
}

func (SideFilterCheckbox) GetTemplateFilename() string {
	return "side_filter_checkbox.gohtml"
}

type SidePriceRange struct {
	Q             string
	FilterParam   anor.FilterParam
	FilteringData anor.FilterParam
}

func (SidePriceRange) GetTemplateFilename() string {
	return "side_price_range.gohtml"
}

type SideRatingCheckbox struct {
}

func (SideRatingCheckbox) GetTemplateFilename() string {
	return "side_rating_checkbox.gohtml"
}

type Pagination struct {
	Q             string
	TotalPages    int
	CurrentPage   int
	TotalProducts int
}

func (Pagination) GetTemplateFilename() string {
	return "pagination.gohtml"
}
