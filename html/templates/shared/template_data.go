package shared

import "github.com/aliml92/anor"

type CategoryBreadcrumb struct {
	Category           anor.Category
	AncestorCategories []anor.Category
}

func (CategoryBreadcrumb) GetTemplateFilename() string {
	return "category_breadcrumb.html"
}

type ProductGrid struct {
	Products []anor.Product
}

func (ProductGrid) GetTemplateFilename() string {
	return "product_grid.gohtml"
}

type Stepper struct {
	CurrentStep int
}

func (Stepper) GetTemplateFilename() string {
	return "stepper.gohtml"
}
