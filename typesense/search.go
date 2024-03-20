package typesense

import (
	"context"
	"fmt"

	"github.com/aliml92/go-typesense/typesense"

	"github.com/aliml92/anor"
)

type QuerySuggestionResults struct {
	Products     []anor.Product
	Categories   []anor.Category
	SellerStores []anor.SellerStore
}

func (s Searcher) SearchQuerySuggestions(ctx context.Context, q string) error {
	// var r QuerySuggestionResults

	searchRequests := &typesense.MultiSearchSearchesParameter{
		Searches: []typesense.MultiSearchCollectionParameters{
			{
				Collection: anor.INDEXPRODUCTS,
				QueryBy:    typesense.String("name"),
				Q:          typesense.String(q),
			},
			{
				Collection: anor.INDEXCATEGORIES,
				QueryBy:    typesense.String("category"),
				Q:          typesense.String(q),
			},
			{
				Collection: anor.INDEXSELLERSTORES,
				QueryBy:    typesense.String("name"),
				Q:          typesense.String(q),
			},
		},
	}

	commonSearchParams := &typesense.MultiSearchParameters{
		Page:    typesense.Int(0),
		PerPage: typesense.Int(5),
	}

	res, err := s.client.Documents.MultiSearch(ctx, searchRequests, commonSearchParams)
	if err != nil {
		return err
	}

	productResults := res.Results[0]
	categoryResults := res.Results[1]
	sellerStoreResults := res.Results[2]

	fmt.Printf("product results: %v\n", productResults.SearchResult.Hits[0].Document)
	fmt.Printf("category results: %v\n", categoryResults.SearchResult)
	fmt.Printf("seller store results: %v\n", sellerStoreResults.SearchResult)

	return nil
}
