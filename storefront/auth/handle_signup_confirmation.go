package auth

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor/html/dtos/pages/signin"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"net/http"
)

type SignupConfirmationForm struct {
	OTP   string
	Email string
}

func (f *SignupConfirmationForm) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	f.OTP = r.PostForm.Get("otp")
	f.Email = r.PostForm.Get(inputEmail)

	return nil
}

func (f *SignupConfirmationForm) Validate() error {
	err := validation.Errors{
		"otp": validation.Validate(f.OTP, validation.Required,
			validation.Length(6, 6),
			is.Digit),
		"email": validation.Validate(f.Email, validation.Required, is.EmailFormat),
	}.Filter()

	return err
}

func (h *Handler) SignupConfirmation(w http.ResponseWriter, r *http.Request) {
	f := &SignupConfirmationForm{}

	err := bindValid(r, f)
	if err != nil {
		h.logClientError(err)
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.svc.SignupConfirm(ctx, f.OTP, f.Email); err != nil {
		switch {
		case errors.Is(err, ErrInvalidOTP):
			err = fmt.Errorf("%s. Please ensure that the OTP is entered correctly and not expired", err.Error())
			h.clientError(w, err, http.StatusBadRequest)
		default:
			h.serverInternalError(w, err)
		}
		return
	}

	message := "You've successfully signed up &#129395"
	sc := signin.Content{Message: formatMessage(message, "success")}
	h.view.Render(w, "pages/signin/content.gohtml", sc)
}
