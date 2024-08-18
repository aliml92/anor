package components

import "github.com/aliml92/anor"

type ProductMain struct {
	Product *anor.Product
}

func (ProductMain) GetTemplateFileName() string {
	return "product_main.gohtml"
}

type ProductSpecs struct {
	Product *anor.Product
}

func (ProductSpecs) GetTemplateFileName() string {
	return "product_specs.gohtml"
}
