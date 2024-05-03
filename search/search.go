package anor

import "context"

const (
	INDEXPRODUCTS       = "products"
	INDEXCATEGORIES     = "categories"
	INDEXSELLERSTORES   = "sellerstores"
	INDEXPRODUCTQUERIES = "product_queries"
)

type SearchQuerySuggestionResults struct {
	ProductSuggestions []string
	Categories         []Category
	SellerStores       []SellerStore
}

type SearchResults struct {
	Products    []Product
	CategoryIDs []int32
	Options     QueryOptions
}

type Searcher interface {
	SearchQuerySuggestions(ctx context.Context, q string) (SearchQuerySuggestionResults, error)
	SearchProducts(ctx context.Context, q string, o QueryOptions) (SearchResults, error)
}
