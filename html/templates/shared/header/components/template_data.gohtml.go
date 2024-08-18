package components

import "html/template"

type WishlistNavItem struct {
	HxSwapOOB          bool
	WishlistItemsCount int
}

func (WishlistNavItem) GetTemplateFilename() string {
	return "wishlist_nav_item.gohtml"
}

type CartNavItem struct {
	HxSwapOOB      bool
	CartItemsCount int
}

func (CartNavItem) GetTemplateFilename() string {
	return "cart_nav_item.gohtml"
}

type SearchQuerySuggestionsList struct {
	ProductNameSuggestions []template.HTML
}

func (SearchQuerySuggestionsList) GetTemplateFilename() string {
	return "search_query_suggestions_list.gohtml"
}
