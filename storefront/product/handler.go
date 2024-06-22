package product

import (
	"fmt"
	"github.com/aliml92/anor/redis/cache/session"
	"github.com/aliml92/anor/search"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
	"github.com/aliml92/anor/pkg/httperrors"
)

var ErrInvalidHandle = errors.New("invalid handle value")

type BindValidator interface {
	Binder
	Validator
}

type Binder interface {
	Bind(r *http.Request) error
}

type Validator interface {
	Validate() error
}

func bindValid[T BindValidator](r *http.Request, v T) error {
	if err := v.Bind(r); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}
	if err := v.Validate(); err != nil {
		return err
	}
	return nil
}

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
