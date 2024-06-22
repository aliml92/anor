package product

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/invopop/validation"
	"github.com/shopspring/decimal"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

const (
	defaultPage         = 1
	defaultPageSize     = 20
	maxFilterableBrands = 10
	maxFilterableColors = 10
	maxFilterableSizes  = 10
)

const (
	sort = "sort"

	filterRating    = "rating"
	filterPriceFrom = "price_from"
	filterPriceTo   = "price_to"
	filterBrands    = "brands"
	filterColors    = "colors"
	filterSizes     = "sizes"

	pagingPage     = "page"
	pagingPageSize = "page_size"
)

type QueryParams struct {
	anor.SortParam
	anor.FilterParam
	anor.Paging
}

func (q *QueryParams) Bind(r *http.Request) error {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return err
	}

	sort := values.Get(sort)
	rating := values.Get(filterRating)
	priceFrom := values.Get(filterPriceFrom)
	priceTo := values.Get(filterPriceTo)
	brands := values.Get(filterBrands)
	colors := values.Get(filterColors)
	sizes := values.Get(filterSizes)
	page := values.Get(pagingPage)
	pageSize := values.Get(pagingPageSize)

	q.SortParam = anor.SortParam(sort)

	var errs error

	if rating != "" {
		ratingInt, err := strconv.Atoi(rating)
		if err != nil {
			err := fmt.Errorf("invalid rating value: '%s'. The rating must be an integer", rating)
			errs = errors.Join(err)
		}
		q.Rating = ratingInt
	}

	if priceFrom != "" {
		priceF, err := strconv.ParseFloat(priceFrom, 64)
		if err != nil {
			err := fmt.Errorf("invalid price_from value: '%s'. The price_from value must be a number", priceFrom)
			errs = errors.Join(err)
		}
		q.PriceFrom = decimal.NewFromFloat(priceF)
	}

	if priceTo != "" {
		priceT, err := strconv.ParseFloat(priceTo, 64)
		if err != nil {
			err := fmt.Errorf("invalid price_to value: '%s'. The price_to value must be a number", priceTo)
			errs = errors.Join(err)
		}
		q.PriceTo = decimal.NewFromFloat(priceT)
	}

	q.Brands = slices.DeleteFunc(strings.Split(brands, ","), func(e string) bool {
		return e == ""
	})
	q.Colors = strings.Split(colors, ",")
	q.Sizes = strings.Split(sizes, ",")

	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			err := fmt.Errorf("invalid page value: '%s'. The page must be an integer", rating)
			errs = errors.Join(err)
		}
		q.Page = pageInt
	} else {
		q.Page = defaultPage
	}

	if pageSize != "" {
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			err := fmt.Errorf("invalid page_size value: '%s'. The page_size must be an integer", rating)
			errs = errors.Join(err)
		}
		q.PageSize = pageSizeInt
	} else {
		q.PageSize = defaultPageSize
	}

	return errs
}

func (q *QueryParams) Validate() error {
	return validation.ValidateStruct(q,
		validation.Field(&q.SortParam,
			validation.In(
				anor.SortParamPopular,
				anor.SortParamPriceHighToLow,
				anor.SortParamPriceLowToHigh,
				anor.SortParamNewArrivals,
				anor.SortParamHighestRated,
				anor.SortParamBestSellers,
			),
		),
		validation.Field(&q.Brands, validation.Length(0, maxFilterableBrands)),
		validation.Field(&q.Colors, validation.Length(0, maxFilterableColors)),
		validation.Field(&q.Sizes, validation.Length(0, maxFilterableSizes)),
		validation.Field(&q.Page, validation.Min(1)),
		validation.Field(&q.PageSize, validation.Min(20), validation.Max(100)),
	)
}
