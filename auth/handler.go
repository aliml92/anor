package auth

import (
	"context"
	"fmt"
	"github.com/aliml92/anor/pkg/httperrors"
	"html/template"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"unicode"

	"github.com/alexedwards/scs/v2"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
)

type Handler struct {
	svc     anor.AuthService
	session *scs.SessionManager
	render  *html.Render
	logger  *slog.Logger
}

func NewHandler(
	svc anor.AuthService,
	templ *html.Render,
	session *scs.SessionManager,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		svc:     svc,
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

	http.Error(w, formatError(err.Error()), statusCode)
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

	http.Error(w, formatError("Something went wrong. Please try again later."), http.StatusInternalServerError)
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

func formatMessage(message string, level string) template.HTML {
	var bsIcon, bsAlertType string
	switch level {
	case "error":
		bsIcon = "x-circle-fill"
		bsAlertType = "danger"
	default:
		bsIcon = "check-circle-fill"
		bsAlertType = "success"
	}
	fm := fmt.Sprintf(`
	  <div class="alert alert-%s d-flex align-items-stretch my-0" role="alert">
        <span class="d-inline-block pt-1">
			<i class="bi bi-%s" style="font-size: 24px;"></i>
		</span>
		<div class="d-inline-block ms-3" style="font-size: 0.875rem">%s</div>
	  </div>
	`, bsAlertType, bsIcon, message)

	return template.HTML(fm)
}

func formatError(errorString string) string {
	fm := fmt.Sprintf(`
	  <div class="alert alert-danger d-flex my-0" role="alert">
        <span class="pt-1">
			<i class="bi bi-x-circle-fill" style="font-size: 24px;"></i>
		</span>
		<div class="ms-3" style="font-size: 0.875rem">%s</div>
	  </div>
	`, errorString)

	return fm
}
