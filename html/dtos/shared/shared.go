package shared

import (
	"github.com/aliml92/anor"
)

type CategoryBreadcrumb struct {
	Category           anor.Category
	AncestorCategories []anor.Category
}

type ProductGrid struct {
	Products []anor.Product
}

type RelatedProducts struct {
}
