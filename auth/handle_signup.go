package auth

import (
	"errors"
	"fmt"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"html/template"
	"net/http"
)

const (
	minAccountNameLength = 1
	maxAccountNameLength = 255
	minPasswordLength    = 8
	maxPasswordLength    = 64
)

const (
	inputName            = "name"
	inputEmail           = "email"
	inputPassword        = "password"
	inputConfirmPassword = "confirm-password"
)

type SignupForm struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}

func (f *SignupForm) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	f.Name = r.PostForm.Get(inputName)
	f.Email = r.PostForm.Get(inputEmail)
	f.Password = r.PostForm.Get(inputPassword)
	f.ConfirmPassword = r.PostForm.Get(inputConfirmPassword)

	return nil
}

func (f *SignupForm) Validate() error {
	err := validation.Errors{
		"name":     validation.Validate(f.Name, validation.Required, validation.Length(minAccountNameLength, maxAccountNameLength)),
		"email":    validation.Validate(f.Email, validation.Required, is.EmailFormat),
		"password": validation.Validate(f.Password, validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
		"confirm password": validation.Validate(f.ConfirmPassword, validation.By(func(value interface{}) error {
			cp := value.(string)
			if cp != f.Password {
				return errors.New("passwords don't match")
			}
			return nil
		})),
	}.Filter()

	return err
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	f := &SignupForm{}

	err := bindValid(r, f)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.svc.Signup(ctx, f.Name, f.Email, f.Password); err != nil {
		if errors.Is(err, ErrEmailAlreadyTaken) {
			err = errors.New("the email address is already registered. Please sign in or use a different email address to signup")
			h.clientError(w, err, http.StatusBadRequest)
			return
		}
		h.serverInternalError(w, err)
		return
	}

	message := fmt.Sprintf("We've sent a one time password (OTP) to %s. If you haven't received the OTP, "+
		"please check your spam folder or request a new one.", f.Email)
	successMessage := formatMessage(message, "success")
	data := struct {
		Email        string
		AlertMessage template.HTML
	}{
		Email:        f.Email,
		AlertMessage: successMessage,
	}
	h.render.HTMX(w, http.StatusAccepted, "signup-confirmation.gohtml", data)
}
