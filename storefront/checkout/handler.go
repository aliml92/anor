package checkout

import (
	"context"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/config"
	"github.com/aliml92/anor/html/templates"
	"github.com/aliml92/anor/redis/session"
	"log/slog"
	"net/http"
	"strings"
	"unicode"

	"github.com/aliml92/anor/html"
)

// HandlerConfig contains the dependencies for creating a new Handler.
type HandlerConfig struct {
	UserService          anor.UserService
	CartService          anor.CartService
	StripePaymentService anor.StripePaymentService
	OrderService         anor.OrderService
	CategorySvc          anor.CategoryService
	AddressService       anor.AddressService
	View                 *html.View
	SessionManager       *session.Manager
	Logger               *slog.Logger
	Config               *config.Config
}

// Handler manages HTTP requests for checkout operations.
type Handler struct {
	userService          anor.UserService
	cartService          anor.CartService
	stripePaymentService anor.StripePaymentService
	orderService         anor.OrderService
	categorySvc          anor.CategoryService
	addressService       anor.AddressService
	session              *session.Manager
	view                 *html.View
	logger               *slog.Logger
	cfg                  *config.Config
}

// NewHandler creates and returns a new Handler instance.
func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		userService:          cfg.UserService,
		cartService:          cfg.CartService,
		stripePaymentService: cfg.StripePaymentService,
		orderService:         cfg.OrderService,
		categorySvc:          cfg.CategorySvc,
		addressService:       cfg.AddressService,
		view:                 cfg.View,
		session:              cfg.SessionManager,
		logger:               cfg.Logger,
		cfg:                  cfg.Config,
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

		base := templates.CheckoutBase{
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

func (h *Handler) jsonClientError(w http.ResponseWriter, err error, statusCode int) {
	anor.JSONClientError(h.logger, w, err, statusCode)
}

func (h *Handler) jsonServerInternalError(w http.ResponseWriter, err error) {
	anor.JSONServerInternalError(h.logger, w, err)
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

func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

func isHXRequest(r *http.Request) bool {
	return r.Header.Get("Hx-Request") == "true"
}
