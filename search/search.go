package search

import (
	"context"
	"errors"
	"github.com/aliml92/anor"
	"github.com/shopspring/decimal"
)

const (
	INDEXPRODUCTS       = "products"
	INDEXCATEGORIES     = "categories"
	INDEXSTORES         = "stores"
	INDEXPRODUCTQUERIES = "product_queries"
)

var ErrNoSearchResults = errors.New("no search results found")

type QuerySuggestionResults struct {
	ProductNameSuggestions []string
	Categories             []anor.Category
	SellerStores           []anor.Store
}

type ProductResults struct {
	Products    []anor.Product
	CategoryIDs []int32
	PriceRange  [2]decimal.Decimal
	Brands      []string
}

type ProductsParams struct {
	Filter anor.FilterParam
	Sort   anor.SortParam
	Paging anor.Paging
}

type Searcher interface {
	SearchQuerySuggestions(ctx context.Context, q string) (QuerySuggestionResults, error)
	SearchProducts(ctx context.Context, q string, params ProductsParams) (ProductResults, int64, error)
}
