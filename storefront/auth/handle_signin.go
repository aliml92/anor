package auth

import (
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/session"
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
	var f SigninForm
	err := anor.BindValid(r, &f)
	if err != nil {
		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u, err := h.authService.Signin(ctx, f.Email, f.Password)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			err = fmt.Errorf("%w. Please check your email and password combination", err)
		case ErrEmailNotConfirmed:
			err = fmt.Errorf("%w. Please verify your email before proceeding. "+
				"<button name='email' class='text-link btn-resend' hx-post='/auth/confirmation/resend' hx-target='#content' value='%s'>Verify Your Email</button>",
				err, f.Email)
		case ErrAccountBlocked:
			err = fmt.Errorf("%w. Contact our support team for assistance", err)
		case ErrOAuth2RegisteredAccount:
			err = fmt.Errorf("this email is associated with a Google Sign-In account. You have two options:\\n1. Sign in with Google\\n2. Set up a password for traditional login in your profile settings after signing in with Google")
			h.clientError(w, err, http.StatusConflict)
			return
		default:
			h.serverInternalError(w, err)
			return
		}

		h.clientError(w, err, http.StatusBadRequest)
		return
	}

	su := session.UserFromContext(ctx)

	var userCart anor.Cart
	guestCartID := su.CartID

	if guestCartID != 0 {
		if userCart, err = h.cartService.Merge(ctx, anor.CartMergeParams{
			GuestCartID: guestCartID,
			UserID:      u.ID,
		}); err != nil {
			h.serverInternalError(w, err)
			return
		}
	} else {
		userCart, err = h.cartService.GetByUserID(ctx, u.ID, false)
		if err != nil && !errors.Is(err, anor.ErrCartNotFound) {
			h.serverInternalError(w, err)
			return
		}
	}

	fmt.Printf("cart ID: %d\n", userCart.ID)
	if err := h.session.SetAuthUser(ctx, session.User{
		ID:        u.ID,
		IsAuth:    true,
		Firstname: u.GetFirstname(),
		CartID:    userCart.ID,
	}); err != nil {
		h.serverInternalError(w, err)
		return
	}

	redirectURL := h.getRedirectURL(r)
	h.redirect(w, redirectURL)
}
