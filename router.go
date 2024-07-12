package anor

import (
	"net/http"
	"strings"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Router struct {
	middlewares []MiddlewareFunc
	basePath    string
	mux         *http.ServeMux
}

func (router *Router) Use(mw ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, mw...)
}

// HandleFunc registers a new route with the given pattern, handler function, and optional middleware functions.
func (router *Router) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	finalHandler := http.Handler(handler)

	finalHandler = router.applyMiddlewares(finalHandler, middlewares...)

	ps := strings.Fields(pattern)
	pattern = ps[0] + " " + router.basePath + ps[1]

	router.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		finalHandler.ServeHTTP(w, r)
	})
}

func (router *Router) Handle(pattern string, handler http.Handler, middlewares ...MiddlewareFunc) {
	finalHandler := handler

	finalHandler = router.applyMiddlewares(finalHandler, middlewares...)

	ps := strings.Fields(pattern)
	pattern = ps[0] + " " + router.basePath + ps[1]

	router.mux.Handle(pattern, finalHandler)
}

// ServeHTTP implements the http.Handler interface.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (router *Router) Group(prefix string, middlewares ...MiddlewareFunc) *Router {
	m := make([]MiddlewareFunc, len(router.middlewares))
	copy(m, router.middlewares)
	r := &Router{
		mux:         router.mux,
		basePath:    router.basePath + prefix,
		middlewares: m,
	}
	r.middlewares = append(r.middlewares, middlewares...)

	return r
}

func (router *Router) applyMiddlewares(handler http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	// Apply per-route middlewares
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	// Apply global middlewares
	for i := len(router.middlewares) - 1; i >= 0; i-- {
		handler = router.middlewares[i](handler)
	}

	return handler
}
