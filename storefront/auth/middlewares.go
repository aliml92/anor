package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/aliml92/anor"
	"net/http"
	"time"
)

func (h *Handler) RedirectAuthUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		previousChainV := r.Context().Value("middlewareChain")
		previousChain, ok := previousChainV.(string)
		if !ok {
			previousChain = "unknown"
		}
		ctx := context.WithValue(r.Context(), "middlewareChain", previousChain+" -> RedirectAuthUserMiddleware")
		r = r.WithContext(ctx)

		path := r.URL.Path

		fmt.Printf("chain: %v, path: %s\n", r.Context().Value("middlewareChain"), path)

		id := h.session.Auth.GetInt64(r.Context(), "authenticatedUserID")
		if id != 0 {
			// retrieve user by id from database
			_, err := h.svc.GetUser(r.Context(), id) // retrieved user
			if errors.Is(err, anor.ErrNotFound) {
				next.ServeHTTP(w, r)
				return
			} else if err != nil {
				h.serverInternalError(w, err)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		previousChainV := r.Context().Value("middlewareChain")
		previousChain, ok := previousChainV.(string)
		if !ok {
			previousChain = "unknown"
		}
		ctx := context.WithValue(r.Context(), "middlewareChain", previousChain+" -> SessionMiddleware")
		r = r.WithContext(ctx)

		path := r.URL.Path

		fmt.Printf("chain: %v, path: %s\n", r.Context().Value("middlewareChain"), path)

		user := anor.UserFromContext(ctx)
		if user == nil {
			var token string
			cookie, err := r.Cookie(h.session.Guest.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			ctx, err = h.session.Guest.Load(r.Context(), token)
			if err != nil {
				h.session.Guest.ErrorFunc(w, r, err)
				return
			}

			if token == "" {
				err := h.session.Guest.RenewToken(ctx)
				if err != nil {
					h.session.Guest.ErrorFunc(w, r, err)
					return
				}
			}

			sr := r.WithContext(ctx)

			sw := &sessionResponseWriter{
				ResponseWriter: w,
				request:        sr,
				authSession:    h.session.Auth,
				anonSession:    h.session.Guest,
			}

			next.ServeHTTP(sw, sr)

			if !sw.written {
				sw.commitAndWriteSessionCookie(w, sr)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

type sessionResponseWriter struct {
	http.ResponseWriter
	request     *http.Request
	authSession *scs.SessionManager
	anonSession *scs.SessionManager
	written     bool
}

// WriteHeader overrides the WriteHeader method of http.ResponseWriter.
func (sw *sessionResponseWriter) WriteHeader(code int) {
	if !sw.written {
		sw.commitAndWriteSessionCookie(sw.ResponseWriter, sw.request)
		sw.written = true
	}
	sw.ResponseWriter.WriteHeader(code)
}

// Write overrides the Write method of http.ResponseWriter.
func (sw *sessionResponseWriter) Write(data []byte) (int, error) {
	if !sw.written {
		sw.commitAndWriteSessionCookie(sw.ResponseWriter, sw.request)
		sw.written = true
	}
	return sw.ResponseWriter.Write(data)
}

func (sw *sessionResponseWriter) commitAndWriteSessionCookie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch sw.anonSession.Status(ctx) {
	case scs.Modified:
		token, expiry, err := sw.anonSession.Commit(ctx)
		if err != nil {
			sw.anonSession.ErrorFunc(w, r, err)
			return
		}

		sw.anonSession.WriteSessionCookie(ctx, w, token, expiry)
	case scs.Destroyed:
		sw.anonSession.WriteSessionCookie(ctx, w, "", time.Time{})
	}
}
