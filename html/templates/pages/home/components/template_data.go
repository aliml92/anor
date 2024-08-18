package components

import "github.com/aliml92/anor"

// Featured represents a collection of FeaturedSelections to be displayed
// prominently on the home page/banner carousels.
type Featured struct {
	Selections []anor.FeaturedSelection
}

func (Featured) GetTemplateFilename() string {
	return "featured.gohtml"
}

// Collection represents a curated group of products displayed on the home page.
// It can represent various types of product groupings such as:
// - Popular products
// - Trending products
// - Recommended products
// - Recently viewed products
// - Any other categorization relevant to the home page or product listings
type Collection struct {
	Title        string
	ResourcePath string
	Products     []anor.Product
}

func (Collection) GetTemplateFilename() string {
	return "collection.gohtml"
}
