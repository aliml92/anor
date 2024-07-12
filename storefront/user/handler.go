package user

import (
	"context"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
	notfoundpage "github.com/aliml92/anor/html/dtos/pages/not_found"
	profilepage "github.com/aliml92/anor/html/dtos/pages/profile"
	"github.com/aliml92/anor/html/dtos/partials"
	"github.com/aliml92/anor/redis/cache/session"
	"log/slog"
	"net/http"
	"path"
)

type Handler struct {
	userSvc anor.UserService
	cartSvc anor.CartService
	session *session.Manager
	view    *html.View
	logger  *slog.Logger
}

func NewHandler(
	userSvc anor.UserService,
	cartSvc anor.CartService,
	view *html.View,
	session *session.Manager,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		userSvc: userSvc,
		cartSvc: cartSvc,
		view:    view,
		session: session,
		logger:  logger,
	}
}
func (h *Handler) Render(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	ctx := r.Context()
	if isHXRequest(r) {
		v := path.Join(page, "content.gohtml")
		h.view.Render(w, v, data)
	}

	hc, err := h.headerContent(ctx)
	if err != nil {
		h.serverInternalError(w, err)
		return
	}

	v := path.Join(page, "base.gohtml")
	switch data.(type) {
	case profilepage.Content:
		c, _ := data.(profilepage.Content)
		base := profilepage.Base{
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

		header.ActiveOrdersCount = ac.ActiveOrdersCount
		header.WishlistItemsCount = ac.WishlistItemsCount
		header.CartItemsCount = ac.CartItemsCount

	} else {
		cartId := h.session.Guest.GetInt64(ctx, "guest_cart_id")
		if cartId != 0 {
			guestCartItemCount, err := h.cartSvc.CountCartItems(ctx, cartId)
			if err != nil {
				return partials.Header{}, err
			}
			header.CartItemsCount = int(guestCartItemCount)
		}
	}

	return header, nil
}

func (h *Handler) clientError(w http.ResponseWriter, err error, statusCode int) {
	anor.ClientError(h.logger, w, err, statusCode)
}

func (h *Handler) serverInternalError(w http.ResponseWriter, err error) {
	anor.ServerInternalError(h.logger, w, err)
}

func (h *Handler) redirect(w http.ResponseWriter, url string) {
	// Log redirection
	h.logger.LogAttrs(
		context.TODO(),
		slog.LevelInfo,
		"redirecting to...",
		slog.String("url", url),
	)

	w.Header().Add("HX-Redirect", url)
	w.WriteHeader(http.StatusOK)
}

func isHXRequest(r *http.Request) bool {
	if r.Header.Get("Hx-Request") == "true" {
		return true
	}

	return false
}
