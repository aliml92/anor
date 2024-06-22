package checkout

import (
	"context"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/config"
	"github.com/aliml92/anor/pkg/httperrors"
	"github.com/aliml92/anor/redis/cache/session"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"unicode"

	"github.com/aliml92/anor/html"
)

type Handler struct {
	userSvc  anor.UserService
	cartSvc  anor.CartService
	orderSvc anor.OrderService
	session  *session.Manager
	view     *html.View
	logger   *slog.Logger
	cfg      *config.Config
}

func NewHandler(
	userSvc anor.UserService,
	cartSvc anor.CartService,
	orderSvc anor.OrderService,
	view *html.View,
	session *session.Manager,
	logger *slog.Logger,
	cfg *config.Config,
) *Handler {
	return &Handler{
		userSvc:  userSvc,
		cartSvc:  cartSvc,
		orderSvc: orderSvc,
		view:     view,
		session:  session,
		logger:   logger,
		cfg:      cfg,
	}
}

func (h *Handler) clientError(w http.ResponseWriter, err error, statusCode int) {
	_, file, no, _ := runtime.Caller(1)
	h.logger.LogAttrs(
		context.TODO(),
		slog.LevelError,
		"client error",
		slog.String("file", file),
		slog.String("line", strconv.Itoa(no)),
		slog.String("status", strconv.Itoa(statusCode)),
		slog.String("error", capitalizeFirst(err.Error())),
	)

	http.Error(w, err.Error(), statusCode)
}

func (h *Handler) serverInternalError(w http.ResponseWriter, err error) {
	_, file, no, _ := runtime.Caller(1)
	h.logger.LogAttrs(
		context.TODO(),
		slog.LevelError,
		"server error",
		slog.String("file", file),
		slog.String("line", strconv.Itoa(no)),
		slog.String("status", strconv.Itoa(http.StatusInternalServerError)),
		slog.String("error", capitalizeFirst(err.Error())),
	)

	http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
}

func (h *Handler) redirect(w http.ResponseWriter, url string) {
	// Log redirection
	h.logger.LogAttrs(
		context.TODO(),
		slog.LevelInfo,
		"redirect",
		slog.String("url", url),
	)

	w.Header().Add("HX-Redirect", url)
	w.WriteHeader(http.StatusOK)
}

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

func (h *Handler) logClientError(err error) {
	httperrors.LogClientError(h.logger, err)
}

func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

func isHXRequest(r *http.Request) bool {
	if r.Header.Get("Hx-Request") == "true" {
		return true
	}

	return false
}
