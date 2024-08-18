package product

import (
	"errors"
	"github.com/aliml92/anor/html"
	productlistings "github.com/aliml92/anor/html/templates/pages/product_listings"
	"github.com/aliml92/anor/html/templates/pages/product_listings/components"
	"github.com/aliml92/anor/html/templates/shared"
	"net/http"
	"strconv"
	"strings"

	"github.com/aliml92/anor"
)

type ListingsViewRequest struct {
	CategoryHandle string      // path param
	Query          QueryParams // common listings query params
}

func (l *ListingsViewRequest) Bind(r *http.Request) error {
	handle := r.PathValue("handle")
	l.CategoryHandle = handle

	q := new(QueryParams)
	if err := q.Bind(r); err != nil {
		return err
	}
	l.Query = *q

	return nil
}

func (l *ListingsViewRequest) Validate() error {
	if err := validateCategoryHandle(l.CategoryHandle); err != nil {
		return errors.Join(err, ErrInvalidHandle)
	}
	if err := l.Query.Validate(); err != nil {
		return err
	}

	return nil
}

// ProductListingsView handles the request for the product details page.
func (h *Handler) ProductListingsView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	target := r.Header.Get("Hx-Target")

	var req ListingsViewRequest
	err := anor.BindValid(r, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidHandle) {
			h.renderNotFound(w, r, "Category not found")
			return
		}
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	categoryID := int32(extractID(req.CategoryHandle))

	// get category
	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		if errors.Is(err, anor.ErrNotFound) {
			h.renderNotFound(w, r, "Category not found")
			return
		}
		h.serverInternalError(w, err)
		return
	}

	getProductsParams := anor.ListByCategoryParams{
		Sort:   req.Query.SortParam,
		Filter: req.Query.FilterParam,
		Paging: req.Query.Paging,
	}

	products, count, err := h.productService.ListByCategory(ctx, category, getProductsParams)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	pg := shared.ProductGrid{Products: products}
	pag := components.Pagination{
		CategoryPath:  formatCategoryPath(category.Handle, category.ID),
		TotalPages:    calcTotalPages(count, req.Query.PageSize),
		CurrentPage:   req.Query.Page,
		TotalProducts: int(count),
	}

	if target == "pagination" || target == "product-grid-with-pagination" {
		comps := []html.Component{
			{"shared/product_grid.gohtml": pg},
			{"pages/product_listings/components/pagination.gohtml": pag},
		}
		h.view.RenderComponents(w, comps)
		return
	}

	ch, err := h.categoryService.GetCategoryHierarchy(ctx, category)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	// TODO: add filter params, price range, color, size, etc
	brands, err := h.productService.ListAllBrandsByCategory(ctx, category)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	// TODO: add filter params, brand, color, size, etc
	minmaxPrice, err := h.productService.GetMinMaxPricesByCategory(ctx, category)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	cb := shared.CategoryBreadcrumb{
		Category:           category,
		AncestorCategories: ch.AncestorCategories,
	}

	scl := components.SideCategoryList{
		Category:           category,
		AncestorCategories: ch.AncestorCategories,
		ChildrenCategories: ch.ChildCategories,
		SiblingCategories:  ch.SiblingCategories,
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
		Pagination:          pag,
	}

	h.Render(w, r, "pages/product_listings/content.gohtml", plc)
}

// extractID extracts id from handle
func extractID(handle string) int64 {
	parts := strings.Split(handle, "-")
	if len(parts) > 1 {
		idStr := parts[len(parts)-1]
		id, _ := strconv.ParseInt(idStr, 10, 64) // Ignore the error
		return id
	}

	id, _ := strconv.ParseInt(handle, 10, 64) // Ignore the error
	return id
}

func formatCategoryPath(handle string, id int32) string {
	return "categories/" + handle + "-" + strconv.Itoa(int(id))
}
