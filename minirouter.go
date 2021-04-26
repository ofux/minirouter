package minirouter

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Middleware wraps an http.Handler, returning a new http.Handler.
type Middleware func(next http.Handler) http.Handler

// Mini adds middlewares on top of httprouter
type Mini struct {
	router *httprouter.Router

	basePath    string
	middlewares []Middleware
}

// New initializes a new Mini.
func New() *Mini {
	r := httprouter.New()
	r.HandleMethodNotAllowed = false
	return &Mini{
		router: r,
	}
}

// Router returns the internal httprouter.Router
func (h *Mini) Router() *httprouter.Router {
	return h.router
}

// SubPath returns a new Mini in which a set of sub-routes can be defined. It can be used for inner
// routes that share a common middleware. It inherits all middlewares and base-path of the parent Mini.
func (h *Mini) SubPath(path string) *Mini {
	var middlewaresCopy []Middleware
	if len(h.middlewares) > 0 {
		middlewaresCopy = make([]Middleware, len(h.middlewares))
		copy(middlewaresCopy, h.middlewares)
	}

	return &Mini{
		router:      h.router,
		basePath:    h.path(path),
		middlewares: middlewaresCopy,
	}
}

// WithMiddleware creates a new child Mini instance with one or more middleware.
func (h *Mini) WithMiddleware(middleware ...Middleware) *Mini {
	newMini := h.SubPath("")
	newMini.middlewares = append(newMini.middlewares, middleware...)
	return newMini
}

// WithHandlerMiddleware creates a new child Mini instance and registers an http.Handler as a middleware.
func (h *Mini) WithHandlerMiddleware(handler http.Handler) *Mini {
	return h.WithMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(w, req)
			next.ServeHTTP(w, req)
		})
	})
}

// Handle registers a handler for the given method and path.
func (h *Mini) Handle(method, path string, handler http.Handler, middleware ...Middleware) {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	for i := len(h.middlewares) - 1; i >= 0; i-- {
		handler = h.middlewares[i](handler)
	}
	h.router.Handler(method, h.path(path), handler)
}

// HandleFunc registers a func handler for the given method and path.
func (h *Mini) HandleFunc(method, path string, handler http.HandlerFunc, middleware ...Middleware) {
	h.Handle(method, path, handler, middleware...)
}

// GET registers a GET handler for the given path.
func (h *Mini) GET(path string, handler http.HandlerFunc, middleware ...Middleware) {
	h.Handle(http.MethodGet, path, handler, middleware...)
}

// PUT registers a PUT handler for the given path.
func (h *Mini) PUT(path string, handler http.HandlerFunc, middleware ...Middleware) {
	h.Handle(http.MethodPut, path, handler, middleware...)
}

// POST registers a POST handler for the given path.
func (h *Mini) POST(path string, handler http.HandlerFunc, middleware ...Middleware) {
	h.Handle(http.MethodPost, path, handler, middleware...)
}

// PATCH registers a PATCH handler for the given path.
func (h *Mini) PATCH(path string, handler http.HandlerFunc, middleware ...Middleware) {
	h.Handle(http.MethodPatch, path, handler, middleware...)
}

// DELETE registers a DELETE handler for the given path.
func (h *Mini) DELETE(path string, handler http.HandlerFunc, middleware ...Middleware) {
	h.Handle(http.MethodDelete, path, handler, middleware...)
}

// OPTIONS registers a OPTIONS handler for the given path.
func (h *Mini) OPTIONS(path string, handler http.HandlerFunc, middleware ...Middleware) {
	h.Handle(http.MethodOptions, path, handler, middleware...)
}

// Params returns the httprouter.Params for req.
// This is just a pass-through to httprouter.ParamsFromContext.
func Params(req *http.Request) httprouter.Params {
	return httprouter.ParamsFromContext(req.Context())
}

func (h *Mini) path(p string) string {
	base := strings.TrimSuffix(h.basePath, "/")

	if p != "" && !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	return base + p
}
