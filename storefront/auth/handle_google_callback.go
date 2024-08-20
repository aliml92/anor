package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/session"
	"github.com/markbates/goth/gothic"
	"net/http"
	"strings"
)

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r = r.WithContext(context.WithValue(ctx, "provider", "google"))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		h.serverInternalError(w, fmt.Errorf("failed to complete user authentication: %w", err))
		return
	}

	u, err := h.authService.SigninWithGoogle(ctx, user.Email, user.Name)
	if err != nil {
		if errors.Is(err, ErrAccountBlocked) {
			err = fmt.Errorf("%w. Contact our support team for assistance", err)
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
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *Handler) getRedirectURL(r *http.Request) string {
	url := r.URL.Query().Get("redirect_url")
	if url != "" {
		return url
	}

	host := r.Host
	if strings.HasSuffix(host, "anor.alisherm.dev") {
		return "https://anor.alisherm.dev"
	}

	return "http://localhost:8008"
}
