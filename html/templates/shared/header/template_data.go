package header

import (
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/shared/header/components"
	"github.com/aliml92/anor/session"
)

type Base struct {
	User            session.User
	RootCategories  []anor.Category
	WishlistNavItem components.WishlistNavItem
	CartNavItem     components.CartNavItem
}

func (Base) GetTemplateFilename() string {
	return "base.gohtml"
}
