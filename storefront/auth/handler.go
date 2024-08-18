package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor/html/templates"
	"github.com/aliml92/anor/redis/session"
	"github.com/fatih/color"
	"github.com/samber/oops"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"unicode"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html"
)

type HandlerConfig struct {
	AuthService anor.AuthService
	CartService anor.CartService
	Session     *session.Manager
	View        *html.View
	Logger      *slog.Logger
}

type Handler struct {
	authService anor.AuthService
	cartService anor.CartService
	session     *session.Manager
	view        *html.View
	logger      *slog.Logger
}

func NewHandler(cfg *HandlerConfig) *Handler {
	return &Handler{
		authService: cfg.AuthService,
		cartService: cfg.CartService,
		session:     cfg.Session,
		view:        cfg.View,
		logger:      cfg.Logger,
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

		base := templates.AuthBase{
			Content: td,
		}

		newTemplatePath := strings.ReplaceAll(templatePath, "content.gohtml", "base.gohtml")
		h.view.Render(w, newTemplatePath, base)
	default:
		h.view.RenderComponent(w, templatePath, td)
	}
}

func (h *Handler) clientError(w http.ResponseWriter, err error, statusCode int) {
	h.logger.Error(
		err.Error(),
		slog.Any("error", err),
	)

	errString := formatError(err).Error()
	http.Error(w, errString, statusCode)
}

func (h *Handler) serverInternalError(w http.ResponseWriter, err error) {
	// TODO: slog not render \n in stacktrace error
	//logger.Error(
	//	err.Error(),
	//	slog.Any("error", err),
	//)
	var oopsErr oops.OopsError
	if ok := errors.As(err, &oopsErr); ok {
		fmt.Println(colorizeStacktrace(oopsErr.Stacktrace()))
	} else {
		fmt.Println(err)
	}

	http.Error(w, formatStr("Something went wrong. Please try again later."), http.StatusInternalServerError)
}

func (h *Handler) logError(err error) {
	anor.LogError(h.logger, err)
}

func (h *Handler) redirect(w http.ResponseWriter, url string) {
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

func formatError(error error) error {
	errorString := capitalizeFirst(error.Error())
	errorString = strings.ReplaceAll(errorString, "\\n", "<br>")
	fm := fmt.Sprintf(`
	  <div class="alert alert-danger d-flex my-0" role="alert">
        <span class="pt-1">
			<i class="bi bi-x-circle-fill" style="font-size: 24px;"></i>
		</span>
		<div class="ms-3" style="font-size: 0.875rem">%s</div>
	  </div>
	`, errorString)

	return errors.New(fm)
}

func formatStr(errorString string) string {
	errorString = capitalizeFirst(errorString)
	errorString = strings.ReplaceAll(errorString, "\\n", "<br>")
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

func colorizeStacktrace(stacktrace string) string {
	lines := strings.Split(stacktrace, "\n")
	colorizedLines := make([]string, 0, len(lines))

	for i, line := range lines {
		if i == 0 {
			// Colorize the main error message
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				colorizedLine := fmt.Sprintf("%s: %s", parts[0], color.RedString(parts[1]))
				colorizedLines = append(colorizedLines, colorizedLine)
			} else {
				colorizedLines = append(colorizedLines, line)
			}
		} else if strings.Contains(line, "at ") {
			// Colorize the file and function information
			parts := strings.SplitN(line, "at ", 2)
			if len(parts) == 2 {
				fileAndFunc := strings.SplitN(parts[1], " ", 2)
				if len(fileAndFunc) == 2 {
					lastSlash := strings.LastIndex(fileAndFunc[0], "/")
					path := fileAndFunc[0][:lastSlash+1]
					fileAndLine := fileAndFunc[0][lastSlash+1:]

					colorizedLine := fmt.Sprintf("%s at %s%s %s",
						parts[0],
						path,
						color.CyanString(fileAndLine),
						color.GreenString(fileAndFunc[1]))
					colorizedLines = append(colorizedLines, colorizedLine)
				} else {
					colorizedLines = append(colorizedLines, line)
				}
			} else {
				colorizedLines = append(colorizedLines, line)
			}
		} else {
			colorizedLines = append(colorizedLines, line)
		}
	}

	return strings.Join(colorizedLines, "\n")
}
