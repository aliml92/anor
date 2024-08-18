package auth

import (
	"context"
	"github.com/markbates/goth/gothic"
	"net/http"
)

func (h *Handler) GoogleSignin(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "provider", "google"))
	gothic.BeginAuthHandler(w, r)
}
