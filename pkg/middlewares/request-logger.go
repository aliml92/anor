package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type RequestLogger struct {
	logger        *slog.Logger
	excludedPaths []string
}

func NewRequestLogger(logger *slog.Logger, excludedPaths []string) *RequestLogger {
	return &RequestLogger{
		logger:        logger,
		excludedPaths: excludedPaths,
	}
}

func (rl *RequestLogger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for _, path := range rl.excludedPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		previousChainV := r.Context().Value("middlewareChain")
		previousChain, ok := previousChainV.(string)
		if !ok {
			previousChain = "unknown"
		}
		ctx := context.WithValue(r.Context(), "middlewareChain", previousChain+" -> RequestLogger")
		r = r.WithContext(ctx)

		path := r.URL.Path

		fmt.Printf("chain: %v, path: %s\n", r.Context().Value("middlewareChain"), path)

		startTime := time.Now()

		// Capture response status code and messages using a custom ResponseWriter
		lrw := newLoggingResponseWriter(w)

		// Call the next handler in the chain
		next.ServeHTTP(lrw, r)

		// Calculate latency time
		latency := time.Since(startTime)

		// Log request details including latency time, response status code, and messages
		rl.logger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"Received_Request",
			slog.String("request_id", r.Header.Get("X-Request-ID")),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			// slog.String("remote_addr", r.RemoteAddr),
			// slog.String("client_ip", getClientIP(r)),
			// slog.String("user_agent", r.UserAgent()),
			slog.Duration("latency", latency),
			slog.Int("status_code", lrw.statusCode),
			slog.Int("body_size", lrw.bodySize),
			slog.String("error", lrw.errorMessage),
		)
	})
}

// loggingResponseWriter is a custom http.ResponseWriter that captures the status code and response message.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bodySize     int    // Size of the response body
	errorMessage string // Error message associated with the response
}

// newLoggingResponseWriter creates a new instance of loggingResponseWriter.
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0, ""}
}

// WriteHeader overrides the WriteHeader method of http.ResponseWriter.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write overrides the Write method of http.ResponseWriter.
func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(data)
	lrw.bodySize += size
	if lrw.statusCode >= 400 && lrw.statusCode < 600 {
		lrw.errorMessage = string(data) // Set errorMessage to the provided error message
	}
	return size, err
}

// getClientIP extracts the client's IP address from the request.
func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
	}
	return ip
}
