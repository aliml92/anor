package user

import (
	"context"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/config"
	"github.com/aliml92/anor/html"
	"github.com/aliml92/anor/html/templates"
	"github.com/aliml92/anor/html/templates/shared/header"
	"github.com/aliml92/anor/redis/session"
	"log/slog"
	"net/http"
	"strings"
)

type HandlerConfig struct {
	UserService       anor.UserService
	OrderService      anor.OrderService
	AddressService    anor.AddressService
	Session           *session.Manager
	View              *html.View
	Logger            *slog.Logger
	Config            *config.Config
	GetHeaderDataFunc func(ctx context.Context) (header.Base, error)
}

type Handler struct {
	userService       anor.UserService
	orderService      anor.OrderService
	addressService    anor.AddressService
	cfg               *config.Config
	session           *session.Manager
	view              *html.View
	logger            *slog.Logger
	getHeaderDataFunc func(ctx context.Context) (header.Base, error)
}

func NewHandler(cfg *HandlerConfig) *Handler {
	return &Handler{
		userService:       cfg.UserService,
		orderService:      cfg.OrderService,
		addressService:    cfg.AddressService,
		cfg:               cfg.Config,
		view:              cfg.View,
		session:           cfg.Session,
		logger:            cfg.Logger,
		getHeaderDataFunc: cfg.GetHeaderDataFunc,
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
