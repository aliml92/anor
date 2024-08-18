package components

import "github.com/aliml92/anor"

type Pagination struct {
	CategoryPath  string
	TotalPages    int
	CurrentPage   int
	TotalProducts int
}

func (Pagination) GetTemplateFilename() string {
	return "pagination.gohtml"
}

type ProductListingsInfo struct {
	CategoryPath  string
	TotalProducts int
	SortParam     anor.SortParam
	FilterParam   anor.FilterParam
}

func (ProductListingsInfo) GetTemplateFilename() string {
	return "product_listings_info.gohtml"
}

type SideCategoryList struct {
	Category           anor.Category
	AncestorCategories []anor.Category
	ChildrenCategories []anor.Category
	SiblingCategories  []anor.Category
}

func (SideCategoryList) GetTemplateFilename() string {
	return "side_category_list.gohtml"
}

type SideFilterCheckbox struct {
	FilterName    string
	CategoryPath  string
	FilteringData anor.FilterParam
	FilterParam   anor.FilterParam
}

func (SideFilterCheckbox) GetTemplateFilename() string {
	return "side_filter_checkbox.gohtml"
}

type SidePriceRange struct {
	CategoryPath  string
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
