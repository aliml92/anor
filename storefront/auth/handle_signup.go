package auth

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/templates/pages/auth/signup/components"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
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
	var f SignupForm
	err := anor.BindValid(r, &f)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	_, err = h.authService.Signup(ctx, f.Name, f.Email, f.Password)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyTaken) {
			err = errors.New("the email address is already registered. Please sign in or use a different email address to signup")
			h.clientError(w, err, http.StatusBadRequest)
			return
		}
		if errors.Is(err, ErrOAuth2RegisteredAccount) {
			err = fmt.Errorf("this email is associated with a Google Sign-In account. You have two options:\\n1. Sign in with Google\\n2. Set up a password for traditional login in your profile settings after signing in with Google")
			h.clientError(w, err, http.StatusConflict)
			return
		}
		h.serverInternalError(w, err)
		return
	}

	message := fmt.Sprintf("We've sent a one time password (OTP) to %s. If you haven't received the OTP, "+
		"please check your spam folder or request a new one.", f.Email)

	sc := components.Confirmation{
		Message: formatMessage(message, "success"),
		Email:   f.Email,
	}

	h.Render(w, r, "pages/auth/signup/components/signup_confirmation.gohtml", sc)
}
