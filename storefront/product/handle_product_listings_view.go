package product

import (
	"errors"
	"github.com/aliml92/anor/html"
	productlistings "github.com/aliml92/anor/html/dtos/pages/product-listings"
	"github.com/aliml92/anor/html/dtos/pages/product-listings/components"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/html/dtos/shared"
	"net/http"
	"strconv"

	"github.com/aliml92/anor"
)

type ListingViewRequest struct {
	CategoryHandle string      // path param
	Query          QueryParams // query params
}

func (l *ListingViewRequest) Bind(r *http.Request) error {
	handle := r.PathValue("handle")
	l.CategoryHandle = handle

	q := new(QueryParams)
	if err := q.Bind(r); err != nil {
		return err
	}
	l.Query = *q

	return nil
}

func (l *ListingViewRequest) Validate() error {
	if err := validateCategoryHandle(l.CategoryHandle); err != nil {
		return errors.Join(err, ErrInvalidHandle)
	}
	if err := l.Query.Validate(); err != nil {
		return err
	}

	return nil
}

func (h *Handler) ProductListingView(w http.ResponseWriter, r *http.Request) {
	req := &ListingViewRequest{}
	err := bindValid(r, req)
	if err != nil {
		h.logClientError(err)
		if errors.Is(err, ErrInvalidHandle) {
			h.viewRenderNotFound(w, r, "Category not found")
			return
		}
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	categoryID := int32(extractID(req.CategoryHandle))

	// get category
	category, err := h.categorySvc.GetCategory(ctx, categoryID)
	if err != nil {
		if errors.Is(err, anor.ErrNotFound) {
			h.viewRenderNotFound(w, r, "Category not found")
			return
		}
		h.serverInternalError(w, err)
		return
	}

	getProductsParams := anor.GetProductsByCategoryParams{
		Sort:   req.Query.SortParam,
		Filter: req.Query.FilterParam,
		Paging: req.Query.Paging,
	}

	products, count, err := h.productSvc.GetProductsByCategory(ctx, category, getProductsParams)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	target := r.Header.Get("Hx-Target")

	pg := shared.ProductGrid{Products: products}
	p := components.Pagination{
		CategoryPath:  formatCategoryPath(category.Handle, category.ID),
		TotalPages:    calcTotalPages(count, req.Query.PageSize),
		CurrentPage:   req.Query.Page,
		TotalProducts: int(count),
	}

	if target == "pagination" || target == "product-grid-with-pagination" {
		comps := []html.Component{
			{"shared/product-grid.gohtml": pg},
			{"pages/product-listings/components/pagination.gohtml": p},
		}
		h.view.RenderComponents(w, comps)
		return
	}

	hierarchy, err := h.categorySvc.GetCategoryHierarchy(ctx, category)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	// TODO: add filter params, price range, color, size, etc
	brands, err := h.productSvc.GetProductBrandsByCategory(ctx, category)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	// TODO: add filter params, brand, color, size, etc
	minmaxPrice, err := h.productSvc.GetMinMaxPricesByCategory(ctx, category)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	cb := shared.CategoryBreadcrumb{
		Category:           category,
		AncestorCategories: hierarchy.AncestorCategories,
	}

	scl := components.SideCategoryList{
		Category:           category,
		AncestorCategories: hierarchy.AncestorCategories,
		ChildrenCategories: hierarchy.ChildCategories,
		SiblingCategories:  hierarchy.SiblingCategories,
	}

	filteringData := anor.FilterParam{
		PriceFrom: minmaxPrice[0],
		PriceTo:   minmaxPrice[1],
		Brands:    brands,
	}

	spr := components.SidePriceRange{
		CategoryPath:  formatCategoryPath(category.Handle, category.ID),
		FilterParam:   req.Query.FilterParam,
		FilteringData: filteringData,
	}

	brandsSfc := components.SideFilterCheckbox{
		FilterName:    "Brands",
		CategoryPath:  formatCategoryPath(category.Handle, category.ID),
		FilteringData: filteringData,
		FilterParam:   req.Query.FilterParam,
	}

	pli := components.ProductListingsInfo{
		CategoryPath:  formatCategoryPath(category.Handle, category.ID),
		TotalProducts: int(count),
		SortParam:     req.Query.SortParam,
		FilterParam:   req.Query.FilterParam,
	}

	plc := productlistings.Content{
		CategoryBreadcrumb:  cb,
		SideCategoryList:    scl,
		SidePriceRange:      spr,
		SideBrandsCheckbox:  brandsSfc,
		SideRatingCheckbox:  components.SideRatingCheckbox{}, // TODO: get ratings
		ProductListingsInfo: pli,
		ProductGrid:         pg,
		Pagination:          p,
	}

	if target == "content" {
		h.view.Render(w, "pages/product-listings/content.gohtml", plc)
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

	plb := productlistings.Base{
		Header:  headerContent,
		Content: plc,
	}

	h.view.Render(w, "pages/product-listings/base.gohtml", plb)
}

func formatCategoryPath(handle string, id int32) string {
	return "categories/" + handle + "-" + strconv.Itoa(int(id))
}
