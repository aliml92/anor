package auth

import (
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/pages/auth/signin"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"net/http"
)

type SendResetPasswordLinkForm struct {
	Email string
}

func (f *SendResetPasswordLinkForm) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	f.Email = r.PostForm.Get(inputEmail)
	return nil
}

func (f *SendResetPasswordLinkForm) Validate() error {
	err := validation.Errors{
		"email": validation.Validate(f.Email, validation.Required, is.EmailFormat),
	}.Filter()

	return err
}

func (h *Handler) SendResetPasswordLink(w http.ResponseWriter, r *http.Request) {
	var f SendResetPasswordLinkForm
	err := anor.BindValid(r, &f)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.authService.SendResetPasswordLink(ctx, f.Email); err != nil {
		switch {
		case errors.Is(err, anor.ErrUserNotFound):
			h.logError(err)
			message := "if this email exists in our server, you will receive a reset password link"
			sc := signin.Content{Message: formatMessage(message, "success")}
			h.view.Render(w, "pages/auth/signin/content.gohtml", sc)
		default:
			h.serverInternalError(w, err)
		}
		return
	}

	message := "if this email exists in our server, you will receive a reset password link"
	sc := signin.Content{Message: formatMessage(message, "success")}
	h.Render(w, r, "pages/auth/signin/content.gohtml", sc)
}
