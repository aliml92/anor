package typesense

import (
	"context"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/go-typesense/typesense"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	queryRegex = regexp.MustCompile(`<mark>([^<]+)<\/mark>(\w*)`)
)

func (ts TsSearcher) SearchQuerySuggestions(ctx context.Context, q string) (anor.SearchQuerySuggestionResults, error) {
	var r anor.SearchQuerySuggestionResults

	searchRequests := &typesense.MultiSearchSearchesParameter{
		Searches: []typesense.MultiSearchCollectionParameters{
			{
				Collection:              anor.INDEXPRODUCTS,
				QueryBy:                 typesense.String("name"),
				Q:                       typesense.String(q),
				Infix:                   typesense.String("fallback"),
				HighlightAffixNumTokens: typesense.Int(7),
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
		PerPage: typesense.Int(10),
	}

	res, err := ts.client.Documents.MultiSearch(ctx, searchRequests, commonSearchParams)
	if err != nil {
		return r, err
	}

	// product search results
	pres := res.Results[0]

	// TODO: return category and seller stores as search query suggestions too
	//cres := res.Results[1]
	//sres := res.Results[2]

	r.ProductSuggestions = ExtractProductQuerySuggestions(pres, q)

	return r, nil
}

type ProductDocument struct {
	ID                    int64             `mapstructure:"id"`
	CategoryID            int32             `mapstructure:"category_id"`
	Name                  string            `mapstructure:"name"`
	Brand                 string            `mapstructure:"brand"`
	Slug                  string            `mapstructure:"slug"`
	ThumbImgUrls          map[string]string `mapstructure:"thumb_img_urls"`
	BasePrice             decimal.Decimal   `mapstructure:"base_price"`
	PriceDiscountedAmount decimal.Decimal   `mapstructure:"price_discounted_amount"`
	Rating                float32           `mapstructure:"rating"`
	NumReviews            int               `mapstructure:"num_reviews"`
	CreatedAt             time.Time         `mapstructure:"created_at"`
	UpdatedAt             time.Time         `mapstructure:"updated_at"`
}

func (ts TsSearcher) SearchProducts(ctx context.Context, q string, o anor.QueryOptions) (anor.SearchResults, error) {
	r := anor.SearchResults{}

	params := &typesense.SearchParameters{
		Q:       q,
		QueryBy: "name",
		Infix:   typesense.String("fallback"),
		Page:    typesense.Int(o.Pagination.Page),
		PerPage: typesense.Int(o.Pagination.PageSize),
	}
	res, err := ts.client.Documents.Search(ctx, "products", params)
	if err != nil {
		return r, err
	}

	for _, hit := range res.Hits {
		doc := hit.Document
		var pd ProductDocument
		if err := Decode(doc, &pd); err != nil {
			return r, err
		}

		r.Products = append(r.Products, anor.Product{
			BaseProduct: anor.BaseProduct{
				ID:           pd.ID,
				CategoryID:   pd.CategoryID,
				Name:         pd.Name,
				Brand:        pd.Brand,
				Slug:         pd.Slug,
				ThumbImgUrls: pd.ThumbImgUrls,
				CreatedAt:    pd.CreatedAt,
				UpdatedAt:    pd.UpdatedAt,
			},
			Rating: anor.Rating{
				Rating:      pd.Rating,
				RatingCount: pd.NumReviews,
			},
			Pricing: anor.ProductPricing{
				BasePrice:        pd.BasePrice,
				DiscountedAmount: pd.PriceDiscountedAmount,
			},
		})
	}

	r.Options.Page = o.Pagination.Page
	r.Options.PageSize = o.Pagination.PageSize
	r.Options.TotalItems = int64(*res.Found)
	r.Options.TotalPages = calcTotalPages(int64(*res.Found), o.Pagination.PageSize)

	return r, nil
}

func toTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

func toDecimalHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(decimal.Decimal{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.Float64:
			return decimal.NewFromFloat(data.(float64)), nil
		default:
			return data, nil
		}
	}
}

func Decode(input map[string]interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			toTimeHookFunc(),
			toDecimalHookFunc(),
		),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}

func calcTotalPages(total int64, perPage int) int {
	a := total / int64(perPage)
	b := total % int64(perPage)
	if b != 0 {
		a++
	}
	return int(a)
}

func ExtractProductQuerySuggestions(res typesense.ResultOrError, q string) []string {

	if res.SearchResult != nil {
		var suggestions [][]string

		hits := res.SearchResult.Hits

		for _, hit := range hits {
			if hit.Highlights != nil {
				if len(hit.Highlights) == 0 {
					continue
				}
				snippet := hit.Highlights[0].Snippet

				// if snippet exists extract query suggestions
				if snippet != nil {
					matches := queryRegex.FindAllStringSubmatch(*snippet, -1)
					var extracted []string
					for _, match := range matches {
						extracted = append(extracted, match[1]+match[2])
					}
					suggestions = append(suggestions, extracted)
				}
			}
		}
		uniqueSuggestions := removeDuplicateSuggestions(suggestions, q)
		return flattenSuggestions(uniqueSuggestions)
	} else {
		return []string{}
	}
}

func removeDuplicateSuggestions(suggestions [][]string, query string) [][]string {
	var unique [][]string

	// A map to track unique inner slices. The key is the string representation of the slice.
	seen := make(map[string]bool)

	for _, suggestion := range suggestions {
		// Deduplicate within the suggestion slice
		dedupedSuggestion := removeDuplicates(suggestion)

		// Reorder words if necessary based on the original query
		if strings.Index(query, " ") != -1 {
			dedupedSuggestion = reorderWords(dedupedSuggestion, query)
		}

		// Convert slice to string to use as a map key
		suggestionKey := fmt.Sprintf("%v", dedupedSuggestion)

		// Check if we've already seen this suggestion
		if !seen[suggestionKey] {
			unique = append(unique, dedupedSuggestion)
			seen[suggestionKey] = true
		}
	}

	return unique
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var unique []string

	// Sort the slice to ensure consistent order
	sort.Strings(slice)

	for _, item := range slice {
		if _, found := seen[item]; !found {
			seen[item] = true
			unique = append(unique, strings.ToLower(item))
		}
	}

	return unique
}

func flattenSuggestions(uniqueProductSuggestions [][]string) []string {
	var flattened []string

	for _, suggestion := range uniqueProductSuggestions {
		// Join the suggestion slice with a space
		flattened = append(flattened, strings.Join(suggestion, " "))
	}

	return flattened
}

func reorderWords(suggestion []string, q string) []string {
	queryWords := strings.Fields(q)

	// Check if both the query and suggestion have 2 words and their first letters switch
	if len(queryWords) == 2 && len(suggestion) == 2 &&
		string(queryWords[0][0]) == string(suggestion[1][0]) &&
		string(queryWords[1][0]) == string(suggestion[0][0]) {
		// Reorder query's words to match the original order
		return []string{suggestion[1], suggestion[0]}
	}

	return suggestion
}
