package middlewares

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

// RequestIDHeader is an HTTP Header which contains the request id.
const RequestIDHeader = "X-Request-Id"

// RequestID is a middleware that ensures each request has a unique identifier.
// It checks for an existing ID in the request header, generates a new UUID if
// absent, and sets the ID in the request context and both request and response
// headers. The ID is accessible via the RequestIDHeader constant in the request
// context and headers.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		r.Header.Set(RequestIDHeader, requestID)
		w.Header().Set(RequestIDHeader, requestID)
		ctx := context.WithValue(r.Context(), RequestIDHeader, requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
