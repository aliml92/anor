package auth

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"html/template"
	"net/http"
)

type ResendOTPForm struct {
	Email string
}

func (f *ResendOTPForm) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	f.Email = r.PostForm.Get(inputEmail)
	return nil
}

func (f *ResendOTPForm) Validate() error {
	err := validation.Errors{
		"email": validation.Validate(f.Email, validation.Required, is.EmailFormat),
	}.Filter()

	return err
}

func (h *Handler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	f := &ResendOTPForm{}

	err := bindValid(r, f)
	if err != nil {
		h.logClientError(err)
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.svc.ResendOTP(ctx, f.Email); err != nil {
		switch {
		case errors.Is(err, anor.ErrNotFound):
			err = errors.New("if this email exists in our server, you will receive a new OTP code")
			h.clientError(w, err, http.StatusBadRequest)
		default:
			h.serverInternalError(w, err)
		}
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
