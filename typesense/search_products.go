package typesense

import (
	"context"
	"fmt"
	"github.com/aliml92/anor/search"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/aliml92/anor"
	"github.com/aliml92/go-typesense/typesense"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ProductDocument struct {
	ID              int64           `mapstructure:"id"`
	CategoryID      int32           `mapstructure:"category_id"`
	Name            string          `mapstructure:"name"`
	Brand           string          `mapstructure:"brand"`
	Handle          string          `mapstructure:"handle"`
	ImageUrls       map[int]string  `mapstructure:"image_urls"`
	BasePrice       decimal.Decimal `mapstructure:"base_price"`
	Discount        decimal.Decimal `mapstructure:"discount"`
	DiscountedPrice decimal.Decimal `mapstructure:"discounted_price"`
	Rating          float32         `mapstructure:"rating"`
	NumReviews      int             `mapstructure:"num_reviews"`
	CreatedAt       time.Time       `mapstructure:"created_at"`
	UpdatedAt       time.Time       `mapstructure:"updated_at"`
}

func (ts Searcher) SearchProducts(ctx context.Context, q string, p search.ProductsParams) (search.ProductResults, int64, error) {
	r := search.ProductResults{}
	var found int64

	page := p.Paging.Page
	pageSize := p.Paging.PageSize

	// search p
	params := &typesense.SearchParameters{
		Q:       q,
		QueryBy: "name",
		//Infix:   typesense.String("fallback"), //TODO: Could not find `name` in the infix index. Make sure to enable infix search by specifying `infix: true` in the schema.
		FacetBy: typesense.String("category_id,brand,discounted_price"),
		Page:    typesense.Int(page),
		PerPage: typesense.Int(pageSize),
	}

	var filterBuilder strings.Builder
	if !p.Filter.PriceFrom.IsZero() && !p.Filter.PriceTo.IsZero() {
		filterBuilder.WriteString(fmt.Sprintf("discounted_price:[%v..%v]", p.Filter.PriceFrom, p.Filter.PriceTo))
	} else if !p.Filter.PriceFrom.IsZero() {
		filterBuilder.WriteString(fmt.Sprintf("discounted_price:>=%v", p.Filter.PriceFrom))
	} else if !p.Filter.PriceTo.IsZero() {
		filterBuilder.WriteString(fmt.Sprintf("discounted_price:<=%v", p.Filter.PriceTo))
	}

	if len(p.Filter.Brands) != 0 {
		if filterBuilder.Len() > 0 {
			filterBuilder.WriteString(" && ")
		}
		filterBuilder.WriteString(fmt.Sprintf("brand:=%v", formatSlice(p.Filter.Brands)))
	}

	params.FilterBy = typesense.String(filterBuilder.String())

	res, err := ts.client.Documents.Search(ctx, search.INDEXPRODUCTS, params)
	if err != nil {
		return r, found, err
	}

	found = int64(*res.Found)
	if found == 0 {
		return r, found, search.ErrNoSearchResults
	}

	//categoryCount := 0
	var categoryIDs []int32
	var minPrice decimal.Decimal
	var maxPrice decimal.Decimal
	var brands []string

	for _, facetHit := range res.FacetCounts {
		f := facetHit.FieldName
		if *f == "category_id" {
			//categoryCount = *facetHit.Stats.TotalValues
			for _, c := range *facetHit.Counts {
				idStr := c.Value
				id, _ := strconv.ParseInt(*idStr, 10, 64)
				categoryIDs = append(categoryIDs, int32(id))
			}
			r.CategoryIDs = categoryIDs
		}
		if *f == "discounted_price" {
			stats := facetHit.Stats
			minPrice = decimal.NewFromFloat32(float32(*stats.Min))
			maxPrice = decimal.NewFromFloat32(float32(*stats.Max))
			r.PriceRange = [2]decimal.Decimal{minPrice, maxPrice}
		}
		if *f == "brand" {
			for _, c := range *facetHit.Counts {
				brands = append(brands, *c.Value)
			}
			r.Brands = brands
		}
	}

	for _, hit := range res.Hits {
		doc := hit.Document
		var pd ProductDocument
		if err := Decode(doc, &pd); err != nil {
			return r, found, err
		}

		r.Products = append(r.Products, anor.Product{
			BaseProduct: anor.BaseProduct{
				ID:         pd.ID,
				CategoryID: pd.CategoryID,
				Name:       pd.Name,
				Brand:      pd.Brand,
				Handle:     pd.Handle,
				ImageUrls:  pd.ImageUrls,
				CreatedAt:  pd.CreatedAt,
				UpdatedAt:  pd.UpdatedAt,
			},
			Rating: anor.Rating{
				Rating:      pd.Rating,
				RatingCount: pd.NumReviews,
			},
			Pricing: anor.ProductPricing{
				BasePrice:       pd.BasePrice,
				Discount:        pd.Discount,
				DiscountedPrice: pd.DiscountedPrice,
			},
		})
	}

	return r, found, nil
}

func toTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
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
		data interface{},
	) (interface{}, error) {
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

func extractBrands(p []anor.Product) []string {
	var (
		uniqueBrands = make(map[string]struct{})
		ub           = "Unbranded"
		ubFound      = false
	)

	for _, v := range p {
		// TODO" remove this when duplicates are handled upon saving them
		b := cases.Title(language.English).String(v.Brand)
		if b != ub {
			uniqueBrands[b] = struct{}{}
		} else {
			ubFound = true
		}

	}

	res := make([]string, len(uniqueBrands))
	i := 0
	for k := range uniqueBrands {
		res[i] = k
		i++
	}

	slices.Sort(res)

	if ubFound {
		res = append(res, ub)
	}

	return res
}

func formatSlice(s []string) string {
	if len(s) == 0 {
		return ""
	}
	var builder strings.Builder
	if len(s) > 1 {
		builder.WriteString("[")
	}
	for i, v := range s {
		if i != 0 {
			builder.WriteString(", ")
		}
		if strings.Contains(v, " ") {
			builder.WriteString("`")
			builder.WriteString(v)
			builder.WriteString("`")
		} else {
			builder.WriteString(v)
		}
	}
	if len(s) > 1 {
		builder.WriteString("]")
	}
	return builder.String()
}
