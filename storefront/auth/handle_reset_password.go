package auth

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/pages/auth/signin"
	"github.com/invopop/validation"
	"net/http"
)

type ResetPasswordForm struct {
	Password        string
	ConfirmPassword string
	Token           string
}

func (f *ResetPasswordForm) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	f.Password = r.PostForm.Get(inputPassword)
	f.ConfirmPassword = r.PostForm.Get(inputConfirmPassword)
	f.Token = r.PostForm.Get("token")

	return nil
}

func (f *ResetPasswordForm) Validate() error {
	err := validation.Errors{
		"password": validation.Validate(f.Password, validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
		"confirm password": validation.Validate(f.ConfirmPassword, validation.By(func(value interface{}) error {
			cp := value.(string)
			if cp != f.Password {
				return errors.New("passwords don't match")
			}
			return nil
		})),
		"token": validation.Validate(f.Token, validation.Required),
	}.Filter()

	return err
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var f ResetPasswordForm
	err := anor.BindValid(r, &f)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.authService.ResetPassword(ctx, f.Token, f.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidOrExpiredResetURL) {
			err := fmt.Errorf("%w. Request a new reset password link to proceed", err)
			h.clientError(w, err, http.StatusBadRequest)
			return
		}

		h.serverInternalError(w, err)
		return
	}

	message := "Your password has been successfully updated! &#129395"
	sc := signin.Content{Message: formatMessage(message, "success")}
	h.Render(w, r, "pages/auth/signin/content.gohtml", sc)
}
