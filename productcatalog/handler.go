package productcatalog

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
	"github.com/aliml92/anor/pkg/httperrors"
	ts "github.com/aliml92/anor/typesense"
)

type Handler struct {
	userSvc     anor.UserService
	productSvc  anor.ProductService
	categorySvc anor.CategoryService
	searcher    *ts.Searcher
	render      *html.Render
	logger      *slog.Logger
	session     *scs.SessionManager
}

func NewHandler(
	userSvc anor.UserService,
	productSvc anor.ProductService,
	categorySvc anor.CategoryService,
	searcher *ts.Searcher,
	render *html.Render,
	logger *slog.Logger,
	session *scs.SessionManager,
) *Handler {
	return &Handler{
		userSvc:     userSvc,
		productSvc:  productSvc,
		categorySvc: categorySvc,
		searcher:    searcher,
		render:      render,
		logger:      logger,
		session:     session,
	}
}

type HomeViewData struct {
	User *anor.User
}

func (h *Handler) HomeView(w http.ResponseWriter, r *http.Request) {
	if isHXRequest(r) {
		h.render.HTMX(w, http.StatusOK, "home.tmpl", nil)
		return
	}

	ctx := r.Context()
	u := ctx.Value("auth").(*anor.User)

	h.render.HTML(w, http.StatusOK, "home.tmpl", HomeViewData{User: u})
}

type ProductDetailsView struct {
	User               *anor.User
	Product            *anor.Product
	Category           anor.Category
	AncestorCategories []anor.Category
}

func (h *Handler) ProductView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	viewData := ProductDetailsView{}
	u := ctx.Value("auth").(*anor.User)
	viewData.User = u

	slug := r.PathValue("slug")
	productID, err := extractProductID(slug)
	if err != nil {
		h.logClientError(err)
		h.render.HTML(w, http.StatusOK, "not-found.tmpl", viewData)
		return
	}

	p, err := h.productSvc.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, anor.ErrNotFound) {
			h.render.HTML(w, http.StatusOK, "not-found.tmpl", viewData)
			return
		}
		h.serverInternalError(w, err)
		return
	}

	c, err := h.categorySvc.GetCategory(ctx, p.CategoryID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	ac, err := h.categorySvc.GetAncestorCategories(ctx, p.CategoryID)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	viewData.Product = p
	viewData.Category = c
	viewData.AncestorCategories = ac

	if isHXRequest(r) {
		w.Header().Add("HX-Trigger-After-Swap", "loadSplide")
		h.render.HTMX(w, http.StatusOK, "product.tmpl", viewData)
		return
	}

	h.render.HTML(w, http.StatusOK, "product.tmpl", viewData)
}

type CategoryView struct {
	User               *anor.User
	Category           *anor.Category
	AncestorCategories []anor.Category
	SiblingCategories  []anor.Category
	ChildCategories    []anor.Category
	Products           []anor.Product
	FilterOptions      anor.FilterOption
	SortOptions        anor.SortOption

	Pagination    Pagination
	ProductsCount int64
}

type Pagination struct {
	Page       int
	PerPage    int
	TotalPages int
}

func (h *Handler) CategoryView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	viewData := CategoryView{}
	u := ctx.Value("auth").(*anor.User)
	viewData.User = u

	slug := r.PathValue("slug")
	categoryID, err := extractCategoryID(slug)
	if err != nil {
		h.logClientError(err)
		h.render.HTML(w, http.StatusOK, "not-found.tmpl", viewData)
		return
	}

	var isReqFromPagination bool
	reqFrom := r.URL.Query().Get("req_from")
	if reqFrom == "btn-pagination" {
		isReqFromPagination = true
	}

	// pagination

	// hardcoded page size
	perPage := 20

	// page number
	pageQuery := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageQuery)
	if page == 0 {
		page++
	}

	offset := (page - 1) * perPage
	limit := perPage

	var (
		category           anor.Category
		ancestorCategories []anor.Category
		siblingCategories  []anor.Category
		childCategories    []anor.Category
		products           []anor.Product
		productsCount      int64
	)

	// get category
	category, err = h.categorySvc.GetCategory(ctx, categoryID)
	if err != nil {
		if errors.Is(err, anor.ErrNotFound) {
			h.render.HTML(w, http.StatusOK, "not-found.tmpl", viewData)
			return
		}
		h.serverInternalError(w, err)
		return
	}

	if !category.IsRoot() && !isReqFromPagination {
		ancestorCategories, err = h.categorySvc.GetAncestorCategories(ctx, categoryID)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}
	}

	if category.IsLeaf() {
		if !isReqFromPagination {
			// get sibling categories
			siblingCategories, err = h.categorySvc.GetSiblingCategories(ctx, categoryID)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}
		}

		// populate params
		params := anor.GetProductsByCategoryParams{
			Offset: offset,
			Limit:  limit,
		}

		products, productsCount, err = h.productSvc.GetProductsByLeafCategoryID(ctx, categoryID, params)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

	} else {
		if !isReqFromPagination {
			childCategories, err = h.categorySvc.GetChildCategories(ctx, category.ID)
			if err != nil {
				h.serverInternalError(w, err)
				return
			}
		}

		// populate params
		params := anor.GetProductsByCategoryParams{
			Offset: offset,
			Limit:  limit,
		}

		products, productsCount, err = h.productSvc.GetProductsByNonLeafCategoryID(ctx, categoryID, params)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

	}

	h.logger.Info("products count:", len(products))

	viewData.Products = products
	viewData.Category = &category
	viewData.AncestorCategories = ancestorCategories
	viewData.ChildCategories = childCategories
	viewData.SiblingCategories = siblingCategories
	viewData.Pagination = Pagination{
		Page:       page,
		PerPage:    perPage,
		TotalPages: calcTotalPages(productsCount, perPage),
	}
	viewData.ProductsCount = productsCount

	if isHXRequest(r) {
		if isReqFromPagination {
			h.render.HTMX(w, http.StatusOK, "category-products-pagination.tmpl", viewData)
			return
		}
		h.render.HTMX(w, http.StatusOK, "category.tmpl", viewData)
		return
	}

	h.render.HTML(w, http.StatusOK, "category.tmpl", viewData)
}

func (h *Handler) SearchQuerySuggestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	searchQuery := r.PostForm.Get("q")
	lastKey := r.PostForm.Get("autocomplete")
	fmt.Printf("1. last key: %v\n", lastKey)
	// TODO: handle the case when q is empty

	if err := h.searcher.SearchQuerySuggestions(ctx, searchQuery); err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, "no error")
}

func (h *Handler) clientError(w http.ResponseWriter, err error, statusCode int) {
	httperrors.ClientError(h.logger, w, err, statusCode)
}

func (h *Handler) serverInternalError(w http.ResponseWriter, err error) {
	httperrors.ServerInternalError(h.logger, w, err)
}

func (h *Handler) logClientError(err error) {
	httperrors.LogClientError(h.logger, err)
}

func calcTotalPages(total int64, perPage int) int {
	a := total / int64(perPage)
	b := total % int64(perPage)
	if b != 0 {
		a++
	}
	return int(a)
}
