package partials

import (
	"github.com/aliml92/anor"
	"html/template"
)

type Header struct {
	User            *anor.User
	RootCategories  []anor.Category
	OrdersNavItem   OrdersNavItem
	WishlistNavItem WishlistNavItem
	CartNavItem     CartNavItem
}

type OrdersNavItem struct {
	HxSwapOOB         string
	ActiveOrdersCount int
}

type WishlistNavItem struct {
	HxSwapOOB          string
	WishlistItemsCount int
}

type CartNavItem struct {
	HxSwapOOB      string
	CartItemsCount int
}

type SearchQuerySuggestionsList struct {
	ProductNameSuggestions []template.HTML
}
