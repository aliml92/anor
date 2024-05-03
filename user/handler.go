package user

import (
	"context"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/pkg/httperrors"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"unicode"

	"github.com/alexedwards/scs/v2"

	"github.com/aliml92/anor/html"
)

type Handler struct {
	userSvc anor.UserService
	session *scs.SessionManager
	render  *html.Render
	logger  *slog.Logger
}

func NewHandler(
	userSvc anor.UserService,
	templ *html.Render,
	session *scs.SessionManager,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		userSvc: userSvc,
		render:  templ,
		session: session,
		logger:  logger,
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
	w.Write([]byte{})
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

func (h *Handler) renderView(w http.ResponseWriter, r *http.Request, status int, page string, data interface{}) {
	if isHXRequest(r) {
		h.render.HTMX(w, status, page, data)
		return
	}

	h.render.HTML(w, status, page, data)
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
