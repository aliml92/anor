package product

import (
	"context"
	"fmt"
	homepage "github.com/aliml92/anor/html/dtos/pages/home"
	notfoundpage "github.com/aliml92/anor/html/dtos/pages/not_found"
	productdetailspage "github.com/aliml92/anor/html/dtos/pages/product_details"
	productlistings "github.com/aliml92/anor/html/dtos/pages/product_listings"
	searchlistings "github.com/aliml92/anor/html/dtos/pages/search_listings"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/redis/cache/session"
	"github.com/aliml92/anor/search"
	"log/slog"
	"net/http"
	"path"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
)

type Handler struct {
	userSvc     anor.UserService
	productSvc  anor.ProductService
	categorySvc anor.CategoryService
	cartSvc     anor.CartService
	searcher    search.Searcher
	view        *html.View
	logger      *slog.Logger
	session     *session.Manager
}

func NewHandler(
	userSvc anor.UserService,
	productSvc anor.ProductService,
	categorySvc anor.CategoryService,
	cartSvc anor.CartService,
	searcher search.Searcher,
	view *html.View,
	logger *slog.Logger,
	session *session.Manager,
) *Handler {
	return &Handler{
		userSvc:     userSvc,
		productSvc:  productSvc,
		categorySvc: categorySvc,
		cartSvc:     cartSvc,
		searcher:    searcher,
		view:        view,
		logger:      logger,
		session:     session,
	}
}

func (h *Handler) Render(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	ctx := r.Context()
	if isHXRequest(r) {
		v := path.Join(page, "content.gohtml")
		h.view.Render(w, v, data)
		return
	}

	hc, err := h.headerContent(ctx)
	if err != nil {
		fmt.Println(err.Error())
		h.serverInternalError(w, err)
		return
	}

	v := path.Join(page, "base.gohtml")
	switch data.(type) {
	case homepage.Content:
		c, _ := data.(homepage.Content)
		base := homepage.Base{
			Header:  hc,
			Content: c,
		}

		h.view.Render(w, v, base)
	case productlistings.Content:
		c, _ := data.(productlistings.Content)
		base := productlistings.Base{
			Header:  hc,
			Content: c,
		}
		h.view.Render(w, v, base)
	case productdetailspage.Content:
		c, ok := data.(productdetailspage.Content)
		if !ok {
			slog.Warn("could not convert the content", "data", data)
		}
		base := productdetailspage.Base{
			Header:  hc,
			Content: c,
		}
		h.view.Render(w, v, base)
	case searchlistings.Content:
		c, _ := data.(searchlistings.Content)
		base := searchlistings.Base{
			Header:  hc,
			Content: c,
		}
		h.view.Render(w, v, base)
	case notfoundpage.Content:
		c, _ := data.(notfoundpage.Content)
		base := notfoundpage.Base{
			Header:  hc,
			Content: c,
		}
		h.view.Render(w, v, base)
	default:
		base := notfoundpage.Base{
			Header: hc,
			Content: notfoundpage.Content{
				Message: "Page not found",
			},
		}
		h.view.Render(w, v, base)
	}
}

func (h *Handler) headerContent(ctx context.Context) (partials.Header, error) {
	u := anor.UserFromContext(ctx)
	header := partials.Header{User: u}
	if u != nil {
		ac, err := h.userSvc.GetUserActivityCounts(ctx, u.ID)
		if err != nil {
			return partials.Header{}, err
		}

		header.CartNavItem = partials.CartNavItem{CartItemsCount: ac.CartItemsCount}
		header.WishlistNavItem = partials.WishlistNavItem{WishlistItemsCount: ac.WishlistItemsCount}
		header.OrdersNavItem = partials.OrdersNavItem{ActiveOrdersCount: ac.ActiveOrdersCount}

	} else {
		cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
		if cartId != 0 {
			guestCartItemCount, err := h.cartSvc.CountCartItems(ctx, cartId)
			if err != nil {
				return partials.Header{}, err
			}
			header.CartNavItem = partials.CartNavItem{CartItemsCount: int(guestCartItemCount)}
		}
	}

	rc, err := h.categorySvc.GetRootCategories(ctx)
	if err != nil {
		return header, err
	}
	header.RootCategories = rc

	return header, nil
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
