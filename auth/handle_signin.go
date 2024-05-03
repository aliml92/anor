package auth

import (
	"fmt"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"net/http"
)

type SigninForm struct {
	Email    string
	Password string
}

func (f *SigninForm) Bind(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	f.Email = r.PostForm.Get(inputEmail)
	f.Password = r.PostForm.Get(inputPassword)
	return nil
}

func (f *SigninForm) Validate() error {
	err := validation.Errors{
		"email":    validation.Validate(f.Email, validation.Required, is.EmailFormat),
		"password": validation.Validate(f.Password, validation.Required),
	}.Filter()

	return err
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	f := &SigninForm{}

	err := bindValid(r, f)
	if err != nil {
		h.logClientError(err)
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID, err := h.svc.Signin(ctx, f.Email, f.Password)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			err = fmt.Errorf("%s. Please check your email and password combination", err.Error())
		case ErrEmailNotConfirmed:
			err = fmt.Errorf("%s. Please verify your email before proceeding. "+
				"<button name='email' class='text-link btn-resend' hx-post='/auth/confirmation/resend' hx-target='#main' value='%s'>Verify Your Email</button>",
				err.Error(), f.Email)
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
