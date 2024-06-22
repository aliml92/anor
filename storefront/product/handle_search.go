package product

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
	searchlistings "github.com/aliml92/anor/html/dtos/pages/search-listings"
	"github.com/aliml92/anor/html/dtos/pages/search-listings/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/html/dtos/shared"
	"github.com/aliml92/anor/search"
	"net/http"
	"net/url"
)

type SearchRequest struct {
	Q     string
	Query QueryParams
}

func (s *SearchRequest) Bind(r *http.Request) error {
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

func (s *SearchRequest) Validate() error {
	if err := validateSearchQuery(s.Q); err != nil {
		return err
	}

	if err := s.Query.Validate(); err != nil {
		return err
	}

	return nil
}

func (h *Handler) SearchView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &SearchRequest{}
	err := bindValid(r, req)
	if err != nil {
		h.logClientError(err)
		if errors.Is(err, ErrInvalidSearchQueryLength) || errors.Is(err, ErrInvalidSearchQueryChars) {
			h.viewRenderNotFound(w, r, "We found nothing")
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
			h.viewRenderNotFound(w, r, "We found nothing")
			return
		}
		h.serverInternalError(w, err)
		return
	}

	target := r.Header.Get("Hx-Target")
	pg := shared.ProductGrid{Products: res.Products}
	p := components.Pagination{
		Q:             req.Q,
		TotalPages:    calcTotalPages(count, req.Query.PageSize),
		CurrentPage:   req.Query.Page,
		TotalProducts: int(count),
	}

	if target == "pagination" || target == "product-grid-with-pagination" {
		comps := []html.Component{
			{"shared/product-grid.gohtml": pg},
			{"pages/search-listings/components/pagination.gohtml": p},
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
	fmt.Printf("brands: %v\n", res.Brands)
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
		Pagination:         p,
	}

	if target == "content" {
		h.view.Render(w, "pages/search-listings/content.gohtml", slc)
		return
	}

	u := anor.UserFromContext(ctx)
	headerContent := partials.Header{User: u}
	if u != nil {
		ac, err := h.userSvc.GetUserActivityCounts(ctx, u.ID)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

		headerContent.ActiveOrdersCount = ac.ActiveOrdersCount
		headerContent.WishlistItemsCount = ac.WishlistItemsCount
		headerContent.CartItemsCount = ac.CartItemsCount

	} else {
		cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
		if cartId != 0 {

			guestCartItemCount, err := h.cartSvc.CountCartItems(ctx, cartId)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}

			headerContent.CartItemsCount = int(guestCartItemCount)
		}
	}

	plb := searchlistings.Base{
		Header:  headerContent,
		Content: slc,
	}

	h.view.Render(w, "pages/search-listings/base.gohtml", plb)
}
