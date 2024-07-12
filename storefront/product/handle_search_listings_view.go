package product

import (
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
	searchlistings "github.com/aliml92/anor/html/dtos/pages/search_listings"
	"github.com/aliml92/anor/html/dtos/pages/search_listings/components"
	"github.com/aliml92/anor/html/dtos/shared"
	"github.com/aliml92/anor/search"
	"log/slog"
	"net/http"
	"net/url"
)

type SearchListingsViewRequest struct {
	Q     string      // search query param
	Query QueryParams // common listings query params
}

func (s *SearchListingsViewRequest) Bind(r *http.Request) error {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return err
	}

	searchQ := values.Get("q")
	s.Q = searchQ

	q := new(QueryParams)
	if err := q.Bind(r); err != nil {
		return err
	}
	s.Query = *q

	return nil
}

func (s *SearchListingsViewRequest) Validate() error {
	if err := validateSearchQuery(s.Q); err != nil {
		return err
	}

	if err := s.Query.Validate(); err != nil {
		return err
	}

	return nil
}

// SearchListingsView handles the search request.
// It returns search listings view page
func (h *Handler) SearchListingsView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	target := r.Header.Get("Hx-Target")

	var req SearchListingsViewRequest
	err := anor.BindValid(r, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidSearchQueryLength) || errors.Is(err, ErrInvalidSearchQueryChars) {
			h.logger.Error(
				err.Error(),
				slog.Any("error", err),
			)
			h.renderNotFound(w, r, "We found nothing")
			return
		}
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	searchParams := search.ProductsParams{
		Filter: req.Query.FilterParam,
		Sort:   req.Query.SortParam,
		Paging: req.Query.Paging,
	}

	res, count, err := h.searcher.SearchProducts(ctx, req.Q, searchParams)
	if err != nil {
		if errors.Is(err, search.ErrNoSearchResults) {
			h.logger.Error(
				err.Error(),
				slog.Any("error", err),
			)
			h.renderNotFound(w, r, "We found nothing")
			return
		}
		h.serverInternalError(w, err)
		return
	}

	pg := shared.ProductGrid{Products: res.Products}
	pag := components.Pagination{
		Q:             req.Q,
		TotalPages:    calcTotalPages(count, req.Query.PageSize),
		CurrentPage:   req.Query.Page,
		TotalProducts: int(count),
	}

	if target == "pagination" || target == "product-grid-with-pagination" {
		comps := []html.Component{
			{"shared/product-grid.gohtml": pg},
			{"pages/search-listings/components/pagination.gohtml": pag},
		}
		h.view.RenderComponents(w, comps)
		return
	}

	cb := components.CategoryBreadcrumb{}
	scl := components.SideCategoryList{}
	filteringData := anor.FilterParam{
		PriceFrom: res.PriceRange[0],
		PriceTo:   res.PriceRange[1],
		Brands:    res.Brands,
	}

	spr := components.SidePriceRange{
		Q:             req.Q,
		FilterParam:   req.Query.FilterParam,
		FilteringData: filteringData,
	}
	brandsSfc := components.SideFilterCheckbox{
		Q:             req.Q,
		FilterName:    "Brands",
		FilteringData: filteringData,
		FilterParam:   req.Query.FilterParam,
	}
	sli := components.SearchListingsInfo{
		Q:               req.Q,
		TotalCategories: len(res.CategoryIDs),
		TotalProducts:   int(count),
		SortParam:       req.Query.SortParam,
		FilterParam:     req.Query.FilterParam,
	}

	if len(res.CategoryIDs) == 1 {
		id := res.CategoryIDs[0]
		category, err := h.categorySvc.GetCategory(ctx, id)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}
		ancestorCategories, err := h.categorySvc.GetAncestorCategories(ctx, id)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}
		siblingCategories, err := h.categorySvc.GetSiblingCategories(ctx, id)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

		// fill in category breadcrumb
		cb.Category = category
		cb.AncestorCategories = ancestorCategories

		// fill in side category list
		scl.Category = category
		scl.AncestorCategories = ancestorCategories
		scl.SiblingCategories = siblingCategories
	} else {
		rootCategories, err := h.categorySvc.GetRootCategories(ctx)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}
		cb.RootCategories = rootCategories
		scl.RootCategories = rootCategories
	}

	slc := searchlistings.Content{
		CategoryBreadcrumb: cb,
		SideCategoryList:   scl,
		SidePriceRange:     spr,
		SideBrandsCheckbox: brandsSfc,
		SideRatingCheckbox: components.SideRatingCheckbox{}, // TODO: get ratings
		SearchListingsInfo: sli,
		ProductGrid:        pg,
		Pagination:         pag,
	}

	h.Render(w, r, "pages/search_listings", slc)
}
