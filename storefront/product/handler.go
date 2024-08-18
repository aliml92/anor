package product

import (
	"context"
	"fmt"
	"github.com/aliml92/anor/html/templates/shared/header"

	//"github.com/aliml92/anor/html/dtos/pages/home"
	//notfound "github.com/aliml92/anor/html/dtos/pages/not_found"
	//productdetails "github.com/aliml92/anor/html/dtos/pages/product_details"
	//productlistings "github.com/aliml92/anor/html/dtos/pages/product_listings"
	//searchlistings "github.com/aliml92/anor/html/dtos/pages/search_listings"
	//"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/html/templates"
	"github.com/aliml92/anor/redis/session"
	"github.com/aliml92/anor/search"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
)

type HandlerConfig struct {
	UserService             anor.UserService
	ProductService          anor.ProductService
	CategoryService         anor.CategoryService
	FeatureSelectionService anor.FeaturedSelectionService
	CartService             anor.CartService
	Session                 *session.Manager
	Searcher                search.Searcher
	View                    *html.View
	Logger                  *slog.Logger
	GetHeaderDataFunc       func(ctx context.Context) (header.Base, error)
}

type Handler struct {
	userService              anor.UserService
	productService           anor.ProductService
	categoryService          anor.CategoryService
	featuredSelectionService anor.FeaturedSelectionService
	cartService              anor.CartService
	session                  *session.Manager
	searcher                 search.Searcher
	view                     *html.View
	logger                   *slog.Logger
	getHeaderDataFunc        func(ctx context.Context) (header.Base, error)
}

func NewHandler(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService:              cfg.UserService,
		productService:           cfg.ProductService,
		categoryService:          cfg.CategoryService,
		featuredSelectionService: cfg.FeatureSelectionService,
		cartService:              cfg.CartService,
		session:                  cfg.Session,
		searcher:                 cfg.Searcher,
		view:                     cfg.View,
		logger:                   cfg.Logger,
		getHeaderDataFunc:        cfg.GetHeaderDataFunc,
	}
}

func (h *Handler) Render(w http.ResponseWriter, r *http.Request, templatePath string, td templates.TemplateData) {
	s := strings.Split(templatePath, "/")
	templateFilename := s[len(s)-1]

	// TODO: remove on production
	if templateFilename != td.GetTemplateFilename() {
		panic(fmt.Sprintf("Template-DTO mismatch: Template '%s' does not match DTO for '%s'",
			templateFilename, td.GetTemplateFilename()))
	}

	switch templateFilename {
	case "base.gohtml":
		h.view.Render(w, templatePath, td)
	case "content.gohtml":
		if isHXRequest(r) {
			h.view.Render(w, templatePath, td)
			return
		}

		ctx := r.Context()
		hc, err := h.getHeaderDataFunc(ctx)
		if err != nil {
			h.serverInternalError(w, err)
			return
		}

		base := templates.Base{
			Header:  hc,
			Content: td,
		}

		newTemplatePath := strings.ReplaceAll(templatePath, "content.gohtml", "base.gohtml")
		h.view.Render(w, newTemplatePath, base)
	default:
		h.view.RenderComponent(w, templatePath, td)
	}
}

func (h *Handler) clientError(w http.ResponseWriter, err error, statusCode int) {
	anor.ClientError(h.logger, w, err, statusCode)
}

func (h *Handler) serverInternalError(w http.ResponseWriter, err error) {
	anor.ServerInternalError(h.logger, w, err)
}

func isHXRequest(r *http.Request) bool {
	if r.Header.Get("Hx-Request") == "true" {
		return true
	}

	return false
}

func calcTotalPages(total int64, perPage int) int {
	a := total / int64(perPage)
	b := total % int64(perPage)
	if b != 0 {
		a++
	}
	return int(a)
}
