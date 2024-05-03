package auth

import (
	"fmt"
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

	message := "You've successfully signed up &#129395"
	h.render.HTMX(w, http.StatusOK, "signin.gohtml", formatMessage(message, "success"))
}
