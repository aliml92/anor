package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"

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

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var f SignupForm

	if err := f.bindAndValidate(r); err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.svc.Signup(ctx, f.Name, f.Email, f.Password); err != nil {
		if errors.Is(err, ErrEmailAlreadyTaken) {
			h.clientError(w, err, http.StatusBadRequest)
		}
		h.serverInternalError(w, err)
		return
	}

	h.render.HTMX(w, http.StatusAccepted, "signup-confirm.tmpl", f.Email)
}

func (h *Handler) SignupView(w http.ResponseWriter, r *http.Request) {
	if isHXRequest(r) {
		h.render.HTMX(w, http.StatusOK, "signup.tmpl", nil)
		return
	}

	h.render.HTML(w, http.StatusOK, "signup.tmpl", nil)
}

func (h *Handler) SignupConfirm(w http.ResponseWriter, r *http.Request) {
	var f SignupConfirmForm

	if err := f.bindAndValidate(r); err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.svc.SignupConfirm(ctx, f.OTP, f.Email); err != nil {
		switch err {
		case ErrInvalidOTP:
			err = fmt.Errorf("%s. Please ensure that the OTP is entered correctly and not expired", err.Error())
			h.clientError(w, err, http.StatusBadRequest)

		case ErrExpiredOTP:
			err = fmt.Errorf("%s. Please request a new OTP", err.Error())
			h.clientError(w, err, http.StatusBadRequest)

		default:
			h.serverInternalError(w, err)
		}

		return
	}

	h.redirect(w, "/signin")
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	var f SigninForm

	if err := f.bindAndValidate(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID, err := h.svc.Signin(ctx, f.Email, f.Password)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			err = fmt.Errorf("%s. Please check your email and password combination", err.Error())
		case ErrEmailNotConfirmed:
			err = fmt.Errorf("%s. Please verify your email before proceeding", err.Error())
		case ErrAccountBlocked:
			err = fmt.Errorf("%s. Contact our support team for assistance", err.Error())
		case ErrAccountInactive:
			err = fmt.Errorf("%s. You can reactivate it by following the instructions in your account settings", err.Error())
		default:
			h.serverInternalError(w, err)
			return
		}

		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.session.RenewToken(ctx); err != nil {
		h.serverInternalError(w, err)
		return
	}

	h.session.Put(ctx, "authenticatedUserID", userID)

	h.redirect(w, "/")
}

func (h *Handler) SigninView(w http.ResponseWriter, r *http.Request) {
	if isHXRequest(r) {
		h.render.HTMX(w, http.StatusOK, "signin.tmpl", nil)
		return
	}

	h.render.HTML(w, http.StatusOK, "signin.tmpl", nil)
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
		slog.String("error", err.Error()),
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
		slog.String("error", err.Error()),
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
