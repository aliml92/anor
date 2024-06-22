package middlewares

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique request ID (UUID) only if it's not already provided
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		previousChainV := r.Context().Value("middlewareChain")
		previousChain, ok := previousChainV.(string)
		if !ok {
			previousChain = "unknown"
		}
		ctx := context.WithValue(r.Context(), "middlewareChain", previousChain+" -> RequestID")
		r = r.WithContext(ctx)

		path := r.URL.Path

		fmt.Printf("chain: %v, path: %s\n", r.Context().Value("middlewareChain"), path)

		// Set the request ID as a request header if it's not already set
		r.Header.Set("X-Request-ID", requestID)

		// Set the same request ID as a response header
		w.Header().Set("X-Request-ID", requestID)
		// Pass the request to the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
