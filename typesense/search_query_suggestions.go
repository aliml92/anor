package typesense

import (
	"context"
	"github.com/aliml92/anor/search"
	"github.com/aliml92/go-typesense/typesense"
)

func (ts Searcher) SearchQuerySuggestions(ctx context.Context, q string) (search.QuerySuggestionResults, error) {
	var r search.QuerySuggestionResults

	searchRequests := &typesense.MultiSearchSearchesParameter{
		Searches: []typesense.MultiSearchCollectionParameters{
			{
				Collection:              search.INDEXPRODUCTQUERIES,
				QueryBy:                 typesense.String("q"),
				Q:                       typesense.String(q),
				Infix:                   typesense.String("fallback"),
				HighlightAffixNumTokens: typesense.Int(7),
				HighlightStartTag:       typesense.String("<strong>"),
				HighlightEndTag:         typesense.String("</strong>"),
			},
			{
				Collection: search.INDEXCATEGORIES,
				QueryBy:    typesense.String("category"),
				Q:          typesense.String(q),
			},
			{
				Collection: search.INDEXSTORES,
				QueryBy:    typesense.String("name"),
				Q:          typesense.String(q),
			},
		},
	}

	commonSearchParams := &typesense.MultiSearchParameters{
		Page:    typesense.Int(0),
		PerPage: typesense.Int(10),
	}

	res, err := ts.client.Documents.MultiSearch(ctx, searchRequests, commonSearchParams)
	if err != nil {
		return r, err
	}

	// product search results
	pres := res.Results[0]

	// TODO: return category and seller stores as search query suggestions too
	// cres := res.ProductResults[1]
	// sres := res.ProductResults[2]

	r.ProductNameSuggestions = ExtractProductQuerySuggestions(pres, q)

	return r, nil
}

func ExtractProductQuerySuggestions(res typesense.ResultOrError, q string) []string {
	if res.SearchResult != nil {
		//var suggestions [][]string
		var snippets []string

		hits := res.SearchResult.Hits

		for _, hit := range hits {
			if hit.Highlights != nil {
				if len(hit.Highlights) == 0 {
					continue
				}
				snippet := hit.Highlights[0].Snippet

				// if snippet exists extract query suggestions
				if snippet != nil {
					snippets = append(snippets, *snippet)
				}
			}
		}
		return snippets
	} else {
		return []string{}
	}
}
