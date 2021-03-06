package minirouter

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Middleware wraps an http.Handler, returning a new http.Handler.
type Middleware func(next http.Handler) http.Handler

// Mini adds middlewares on top of httprouter.Router
type Mini struct {
	router *httprouter.Router

	basePath    string
	middlewares []Middleware
}

// New initializes a new Mini.
func New() *Mini {
	r := httprouter.New()
	r.RedirectFixedPath = false // Disable path auto-correction. Let's be strict.
	return &Mini{
		router: r,
	}
}

// ServeHTTP makes Mini implement the http.Handler interface.
func (m *Mini) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.router.ServeHTTP(w, req)
}

// Router returns the internal httprouter.Router
func (m *Mini) Router() *httprouter.Router {
	return m.router
}

// WithBasePath returns a a copy of parent with an augmented base-path, in which a set of sub-routes can be defined.
// It can be used for inner routes that share a common base-path.
// The new base-path is the concatenation of the parent's base-path and the given path (eg. <parent's base-path>/<path>).
func (m *Mini) WithBasePath(path string) *Mini {
	var middlewaresCopy []Middleware
	if len(m.middlewares) > 0 {
		middlewaresCopy = make([]Middleware, len(m.middlewares))
		copy(middlewaresCopy, m.middlewares)
	}

	return &Mini{
		router:      m.router,
		basePath:    m.path(path),
		middlewares: middlewaresCopy,
	}
}

// WithMiddleware returns a copy of parent with one or more new middlewares. It can be used for
// routes that share common middlewares.
func (m *Mini) WithMiddleware(middleware ...Middleware) *Mini {
	newMini := m.WithBasePath("")
	newMini.middlewares = append(newMini.middlewares, middleware...)
	return newMini
}

// WithHandlerMiddleware returns a copy of parent with an http.Handler as a new middleware.
func (m *Mini) WithHandlerMiddleware(handler http.Handler) *Mini {
	return m.WithMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(w, req)
			next.ServeHTTP(w, req)
		})
	})
}

// Handle registers a handler for the given method and path.
func (m *Mini) Handle(method, path string, handler http.Handler, middleware ...Middleware) {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		handler = m.middlewares[i](handler)
	}
	m.router.Handler(method, m.path(path), handler)
}

// HandleFunc registers a func handler for the given method and path.
func (m *Mini) HandleFunc(method, path string, handler http.HandlerFunc, middleware ...Middleware) {
	m.Handle(method, path, handler, middleware...)
}

// GET registers a GET func handler for the given path.
func (m *Mini) GET(path string, handler http.HandlerFunc, middleware ...Middleware) {
	m.Handle(http.MethodGet, path, handler, middleware...)
}

// PUT registers a PUT func handler for the given path.
func (m *Mini) PUT(path string, handler http.HandlerFunc, middleware ...Middleware) {
	m.Handle(http.MethodPut, path, handler, middleware...)
}

// POST registers a POST func handler for the given path.
func (m *Mini) POST(path string, handler http.HandlerFunc, middleware ...Middleware) {
	m.Handle(http.MethodPost, path, handler, middleware...)
}

// PATCH registers a PATCH func handler for the given path.
func (m *Mini) PATCH(path string, handler http.HandlerFunc, middleware ...Middleware) {
	m.Handle(http.MethodPatch, path, handler, middleware...)
}

// DELETE registers a DELETE func handler for the given path.
func (m *Mini) DELETE(path string, handler http.HandlerFunc, middleware ...Middleware) {
	m.Handle(http.MethodDelete, path, handler, middleware...)
}

// OPTIONS registers a OPTIONS func handler for the given path.
func (m *Mini) OPTIONS(path string, handler http.HandlerFunc, middleware ...Middleware) {
	m.Handle(http.MethodOptions, path, handler, middleware...)
}

// Params returns the httprouter.Params for request.
// This is just a pass-through to httprouter.ParamsFromContext.
func Params(req *http.Request) httprouter.Params {
	return httprouter.ParamsFromContext(req.Context())
}

func (m *Mini) path(p string) string {
	if p == "" {
		return m.basePath
	}

	base := strings.TrimSuffix(m.basePath, "/")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return base + p
}
